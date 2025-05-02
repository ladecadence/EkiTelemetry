#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <inttypes.h>
#include "ssdv.h"
#include "decode.h"

decode_info_t decode_ssdv_file(char* file_in, char* file_out) {
    decode_info_t info;
    int i=0;
    FILE *fin = stdin;
    FILE *fout = stdout;
    int errors;
    ssdv_t ssdv;

    uint8_t pkt[SSDV_PKT_SIZE],  *jpeg;
    size_t jpeg_length;

    //callsign[0] = '\0';

    // decoding
    // open files
    fin = fopen(file_in, "rb");
    if (!fin)
        return info;
    fout = fopen(file_out, "wb");
    if (!fout)
        return info;

    ssdv_dec_init(&ssdv, SSDV_PKT_SIZE);
    jpeg_length = 1024 * 1024 * 4;
    jpeg = (uint8_t*) malloc(jpeg_length);
    ssdv_dec_set_buffer(&ssdv, jpeg, jpeg_length);
    while(fread(pkt, 1, SSDV_PKT_SIZE, fin) > 0)
    {
        /* Test the packet is valid */
        if(ssdv_dec_is_packet(pkt, SSDV_PKT_SIZE, &errors) != 0) continue;

        ssdv_packet_info_t p;

        ssdv_dec_header(&p, pkt);
        // console->append(QString::asprintf("Decoded image packet. Callsign: %s, Image ID: %d, Resolution: %dx%d, Packet ID: %d (%d errors corrected)\n"
        //                     ">> Type: %d, Quality: %d, EOI: %d, MCU Mode: %d, MCU Offset: %d, MCU ID: %d/%d\n",
        //             p.callsign_s,
        //             p.image_id,
        //             p.width,
        //             p.height,
        //             p.packet_id,
        //             errors,
        //             p.type,
        //             p.quality,
        //             p.eoi,
        //             p.mcu_mode,
        //             p.mcu_offset,
        //             p.mcu_id,
        //             p.mcu_count
        //             ));

        info.packet_id = p.packet_id;
        info.image_id = p.image_id;
        /* Feed it to the decoder */
        ssdv_dec_feed(&ssdv, pkt);
        i++;
    }

    ssdv_dec_get_jpeg(&ssdv, &jpeg, &jpeg_length);
    fwrite(jpeg, 1, jpeg_length, fout);
    fclose(fout);
    free(jpeg);
    info.decoded_packets = i;
    return info;

    // ui->ssdvStatusLabel->setText(QString("Decodificados ")
    //                              .append(QString::number(i))
    //                              .append(" paquetes de ").
    //                              append(QFileInfo(path).fileName()));

    // updateImage(path);
}