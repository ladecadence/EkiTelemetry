#ifndef _DECODER_H
#define _DECODER_H 

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <inttypes.h>
#include "ssdv.h"

typedef struct decode_info {
    int decoded_packets;
    int image_id;
    int packet_id;
} decode_info_t;

decode_info_t decode_ssdv_file(char* file_in, char* file_out);

#endif