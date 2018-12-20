import numpy as np
import torch
import torch.nn as nn

class DropBlock(nn.Module):
    def __init__(self, keep_prob, block_size):
        super(DropBlock, self).__init__()
        self.keep_prob = keep_prob
        self.block_size = block_size
        self.gamma = 0.

    def forward(self, x):
        if(self.training==False):
            return x
        gamma = self.calc_gamma(x)
        init_mask = torch.zeros((x.shape[0], x.shape[2], x.shape[3])).add(1. - gamma)
        init_mask = init_mask.to(x.device)
        init_mask = torch.bernoulli(init_mask)
        mask = self.create_block_mask(init_mask, self.block_size//2)
        output = x * mask[:,None,:,:]
        output = output * (output.shape[0] * output.shape[2] * output.shape[3]) / output.sum()
        return output
    def calc_gamma(self, x):
        gamma = ((1. - self.keep_prob) * x.shape[2] ** 2) / (self.block_size ** 2 * (x.shape[2] - self.block_size + 1)**2 )
        return gamma

    def create_block_mask(self, m, kernel_size):
        width = m.shape[1]
        height = m.shape[2]
        out = m.clone()
        for n in range(0, m.shape[0]):
            for i in range(0,width):
                for j in range(0,height):
                    if(m[n, i, j] < 1):
                        out[n, np.maximum(0,i-kernel_size):np.minimum(width,i+kernel_size),
                         np.maximum(0,j-kernel_size):np.minimum(height,j+kernel_size)] = 0
        return out
