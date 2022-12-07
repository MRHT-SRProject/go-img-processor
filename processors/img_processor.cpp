#include "img_processor.h"
#include <opencv2/core/cvstd.hpp>
#include <opencv2/core/mat.hpp>
#include <opencv2/imgcodecs/imgcodecs.hpp>
#include <opencv2/imgproc/imgproc.hpp>
#include <opencv2/photo/photo.hpp>
#include <opencv2/video/tracking.hpp>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <vector>

CPixels cvMatToCPixels(cv::Mat);
Buffer vecToCBuf(std::vector<u_int8_t> &);
CPixels getPixels(CMat mat);

void createExceptionResponse(const std::exception &, Response &);

/**
 * @brief Create a Exception Response object
 * 
 * @param e Exception
 * @param resp Response object to update
 */
void createExceptionResponse(const std::exception &e, Response &resp)
{
    char *strexcept = (char *)std::calloc(sizeof(e.what()), 1);
    strcpy(strexcept, e.what());
    resp.data = (void *)strexcept;
    resp.dataType = "STRING";
    resp.len = sizeof(e.what());
    resp.status = 1;
}

/**
 * @brief Given a black and white image matrix, colorize it according to map
 * 
 * @param img The CMatrix representing the bw image to be colorized
 * @param map The number representing the color map as defined by cv::ColormapTypes
 * See https://docs.opencv.org/4.x/d3/d50/group__imgproc__colormap.html
 * @return Response A Response object with the colorized matrix as the data
 */
Response colorize(CMat img, uint8_t map)
{
    Response resp;

    try
    {
        cv::Mat *colorized = new cv::Mat;
        cv::applyColorMap(*((cv::Mat *)img.mat), *colorized, map);
        resp.data = colorized;
        resp.dataType = "MATRIX";
        resp.len = (size_t)colorized->total(0);
        resp.status = 0;
    }
    catch (const std::exception &e)
    {
        createExceptionResponse(e, resp);
    }

    return resp;
}
/**
 * @brief Given a list of images, stack them using ECC algorithm
 * 
 * @param img List of images to stack
 * @param len the number of images in the stack
 * @return Response Response object with data being the final stacked image data
 */
Response stackImages(Image *img, size_t len)
{
    cv::Mat firstImg;
    cv::Mat stacked;
    Response resp;
    try
    {
        for (int i = 0; i < len; i++)
        {
            // Transformation matrix
            cv::Mat M = cv::Mat::eye(3, 3, CV_32F);
            cv::Mat pixels = cv::Mat(img[i].y1 - img[i].y0, img[i].x1 - img[i].x0, CV_8UC1, img[i].pixels);
            if (!i)
            {
                firstImg = pixels;
                stacked = pixels;
                continue;
            }
            cv::Mat warped;
            // create the transformation matrix
            const int iterations = 100;         // 300 iterations for alignment
            const double terminationEps = 1e-6; // Threshold in correleation
            // define termination criteria
            cv::TermCriteria criteria(cv::TermCriteria::COUNT + cv::TermCriteria::EPS, iterations, terminationEps);
            cv::findTransformECC(firstImg, pixels, M, cv::MOTION_HOMOGRAPHY, criteria);
            // warp the image according ot the transformation matrix
            cv::warpPerspective(pixels, warped, M, pixels.size(), cv::INTER_LINEAR + cv::WARP_INVERSE_MAP);
            // stack the image
            stacked += warped;
        }

        // cv::applyColorMap(stacked, colorized, cv::COLORMAP_PLASMA);
        resp.data = new cv::Mat(stacked.clone());
        resp.dataType = "MATRIX";
        resp.len = (size_t)((((cv::Mat *)(resp.data))->total(0)));
        resp.status = 0;
    }
    catch (const std::exception &e)
    {
        createExceptionResponse(e, resp);
    }

    return resp;
}
/**
 * @brief Get the pixels from a CMat
 * 
 * @param mat The matrix to get the pixels from
 * @return CPixels a buffer of pixels
 */
CPixels getPixels(CMat mat)
{
    cv::Mat cvmat = *(cv::Mat *)(mat.mat);
    return cvMatToCPixels(cvmat);
}

/**
 * @brief Given an image, convert the image to grayscale
 * 
 * @param img The Image object from GO
 * @return Response Response object with the grayscale CMat as the data
 */
Response grayscale(Image img)
{
    Response resp;
    try
    {
        cv::Mat pixels = cv::Mat(img.y1 - img.y0, img.x1 - img.x0, CV_8UC4, img.pixels);
        cv::Mat *gray = new cv::Mat();
        cv::cvtColor(pixels, *gray, cv::COLOR_RGBA2GRAY);
        resp.data = gray;
        resp.dataType = "MATRIX";
        resp.len = gray->total(0);
        resp.status = 0;

    } catch(const std::exception& e) {
        createExceptionResponse(e, resp);
    }

    return resp;
}

/**
 * @brief Convert a CMat object to a buffer of PPM data
 * 
 * @param mat The CMat to convert
 * @return Response A Response object with the data holding a buffer of PPM data
 */
Response cMatToImg(CMat mat)
{
    Response resp;
    try {
    std::vector<uint8_t> buf;
    cv::Mat *m = (cv::Mat *)(mat.mat);
    cv::imencode(".ppm", *(cv::Mat *)(mat.mat), buf);
    resp.data = vecToCBuf(buf).data;
    resp.len = buf.size();
    resp.dataType = "PPM";
    resp.status = 0;
    } catch(const std::exception& e) {
        createExceptionResponse(e, resp);
    }

    return resp;
}

/**
 * @brief Helper function to get pixels from a cv::Mat object
 * 
 * @param mat A cv::Mat object to extract pixel data from
 * @return CPixels 
 */
CPixels cvMatToCPixels(cv::Mat mat)
{
    return (CPixels){
        .pixels = mat.data,
        .len = mat.total(0)};
}

/**
 * @brief Helper function to convert a C++ vector to a C style array
 * 
 * @param vec 
 * @return Buffer 
 */
Buffer vecToCBuf(std::vector<uint8_t> &vec)
{
    uint8_t *cbuf = (uint8_t *)malloc(vec.size());
    for (int i = 0; i < vec.size(); i++)
    {
        cbuf[i] = vec[i];
    }
    return Buffer{
        .data = cbuf,
        .len = vec.size()};
}