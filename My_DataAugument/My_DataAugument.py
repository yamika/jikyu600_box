# -*- coding: utf-8 -*-
from PIL import Image
from PIL import ImageFilter

# 画像の反転
def flip_image(img,flip_type='LR'):
    out_img = img.copy()
    if(flip_type == 'LR'):
        out_img = img.transpose(Image.FLIP_LEFT_RIGHT)
    elif(flip_type == 'TB'):
        out_img = img.transpose(Image.FLIP_TOP_BOTTOM)
    return out_img

# 画像の回転
def rotate_image(img,rotate_value=0):
    out_img = img.copy()
    if(rotate_value > 0):
        out_img = img.rotate(rotate_value)
    return out_img

# 画像をずらす
def move_image(img,move_value=0,move_type='left'):
    out_img = img.copy()
    if(move_value > 0 and move_type=='right'):
        for x in range(0,move_value):
            for y in range(0,img.size[1]):
                r,g,b = img.getpixel((0,y))
                out_img.putpixel((x,y),(r,g,b))
                #out_img[x][y] = img[0][y]

        for x in range(move_value,img.size[0]):
            for y in range(0,img.size[1]):
                r,g,b = img.getpixel((x-move_value,y))
                out_img.putpixel((x,y),(r,g,b))
                #out_img[x][y] = img[x - move_value][y]

    if(move_value > 0 and move_type == 'left'):

        for x in range(0,img.size[0]-move_value):
            for y in range(0,img.size[1]):
                r,g,b = img.getpixel((x+move_value,y))
                out_img.putpixel((x,y),(r,g,b))
                #out_img[x][y] = img[0][y]

        for x in range(img.size[0]-move_value,img.size[0]):
            for y in range(0,img.size[1]):
                r,g,b = out_img.getpixel((img.size[0]-move_value-1,y))
                out_img.putpixel((x,y),(r,g,b))
                #out_img[x][y] = img[x - move_value][y]

    if(move_value > 0 and move_type == 'top'):

        for y in range(0,img.size[1]-move_value):
            for x in range(0,img.size[0]):
                r,g,b = img.getpixel((x,y+move_value))
                out_img.putpixel((x,y),(r,g,b))
                #out_img[x][y] = img[0][y]

        for y in range(img.size[1]-move_value,img.size[1]):
            for x in range(0,img.size[0]):
                r,g,b = out_img.getpixel((x,img.size[1]-move_value-1))
                out_img.putpixel((x,y),(r,g,b))
                #out_img[x][y] = img[x - move_value][y]

    if(move_value > 0 and move_type == 'bottom'):

        for y in range(move_value,img.size[1]):
            for x in range(0,img.size[0]):
                r,g,b = img.getpixel((x,y-move_value))
                out_img.putpixel((x,y),(r,g,b))
                #out_img[x][y] = img[0][y]

        for y in range(0,move_value):
            for x in range(0,img.size[0]):
                r,g,b = out_img.getpixel((x,move_value+1))
                out_img.putpixel((x,y),(r,g,b))
                #out_img[x][y] = img[x - move_value][y]

    return out_img

# 画像のコントラストを変更
def change_contrast(img,contrast_value=1):

    out_img = img.copy()

    for x in range(0,img.size[0]):
        for y in range(0,img.size[1]):
            r,g,b = img.getpixel((x,y))
            r = int(r*contrast_value)
            g = int(g*contrast_value)
            b = int(b*contrast_value)
            out_img.putpixel((x,y),(r,g,b))

    return out_img

# 画像をぼかす
def BLUR_filter(img):
    out_img = img.copy()
    #out_img = out_img.filter(ImageFilter.BLUR)
    for x in range(0,out_img.size[0]):
        for y in range(0,out_img.size[1]):
            r,g,b = img.getpixel((x,y))

            r1 = r2 = r3 = r4 = r5 = r6 = r7 = r8 = r;
            g1 = g2 = g3 = g4 = g5 = g6 = g7 = g8 =  g;
            b1 = b2 = b3 = b4 = b5 = b6 = b7 = b8 = b;

            #近傍座標の値を取得

            #1
            if x-1 > 0 and y+1 < img.size[1]:
                r1,g1,b1 = img.getpixel((x-1,y+1))

            #2
            if y+1 < img.size[1]:
                r2,g2,b2 = img.getpixel((x,y+1))

            #3
            if x+1 < img.size[0] and y+1 < img.size[1]:
                r3,g3,b3 = img.getpixel((x+1,y+1))

            #4
            if x-1 > 0:
                r4,g4,b4 = img.getpixel((x-1,y))

            #5
            if x+1 < img.size[0]:
                r5,g5,b5 = img.getpixel((x+1,y))

            #6
            if x-1 > 0 and y-1 > 0:
                r6,g6,b6 = img.getpixel((x-1,y-1))

            #7
            if y-1 > 0:
                r7,g7,b7 = img.getpixel((x,y-1))

            #8
            if x+1 < img.size[0] and y-1 > 0:
                r8,g8,b8 = img.getpixel((x+1,y-1))


            #近傍のRGBを平均化
            r = (r + r1 + r2 + r3 + r4 + r5 + r6 + r7 + r8)/9
            g = (g + g1 + g2 + g3 + g4 + g5 + g6 + g7 + g8)/9
            b = (b + b1 + b2 + b3 + b4 + b5 + b6 + b7 + b8)/9

            out_img.putpixel((x,y),(int(r),int(g),int(b)))

    return out_img
