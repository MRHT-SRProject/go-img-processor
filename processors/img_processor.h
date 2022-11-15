#ifndef __STACK__
#define __STACK__


#ifdef __cplusplus
extern "C"
{
#endif
#include <stdlib.h>
#include <inttypes.h>
    typedef struct CMat {
        void* mat;
    } CMat;

    typedef struct CPixels {
        unsigned char* pixels;
        size_t len;
    } CPixels;

    typedef struct Image
    {
        uint8_t *pixels;
        float exposure; 
        uint32_t x0;
        uint32_t x1;
        uint32_t y0;
        uint32_t y1;
    } Image;

    void stackImages(Image *img, size_t len);
    CMat grayscale(Image img);
    CPixels getPixels(CMat mat);

#ifdef __cplusplus
}
#endif
#endif
