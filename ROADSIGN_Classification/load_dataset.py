
# coding: utf-8

from PIL import Image
import numpy as np
import os
# coding: utf-8
import os.path

def _load_label(file_dir):
    data_files = os.listdir(file_dir)
    labels = []
    for file in data_files:
        l = file.split('-')[0]
        if(l == 'stop'):
            labels.append(0)
        elif(l == 'limitspeed'):
            labels.append(1)
            
    print("Load label : Done!")
    
    return np.array(labels)

def _load_img(file_dir,convert_type='L'):
    data_files = os.listdir(file_dir)
    imgs = []
    for file in data_files:
        img = np.frombuffer(np.array(Image.open(file_dir+'/'+file).convert(convert_type)),dtype=np.uint8)
        imgs.append(img)
    print("Load img : Done!")
    
    return np.array(imgs)

def _change_one_hot_label(X):
    T = np.zeros((X.size, 2))
    for idx, row in enumerate(T):
        row[X[idx]] = 1
        
    return T

def load_dataset(DIR_PATH,convert_type='L',normalize=True,flatten=True,one_hot_label=False):
    labels = _load_label(DIR_PATH)
    imgs = _load_img(DIR_PATH,convert_type)
    
    if normalize:
        imgs = imgs.astype(np.float32)
        imgs /= 255.0
            
    if not flatten:
        if(convert_type == 'L'):
            imgs = imgs.reshape(-1,1,64,64)
        elif(convert_type == 'RGB'):
            imgs = imgs.reshape(-1,3,64,64)
            
                    
    if one_hot_label:
        labels = _change_one_hot_label(labels)        
        
    return imgs,labels



