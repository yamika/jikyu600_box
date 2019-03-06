import torch
import torch.nn as nn
import torch.nn.functional as F

class MultiResBlock(nn.Module):

    def __init__(self, in_channels, out_channels, alpha = 1.67):
        super(MultiResBlock, self).__init__()

        W = alpha * out_channels

        self.conv3x3 = nn.Conv2d(in_channels, out_channels = int(W*0.167), kernel_size=3, stride=1, padding=1, bias=False)
        self.bn3x3 = nn.BatchNorm2d(int(W*0.167))

        self.conv5x5 = nn.Conv2d(int(W*0.167), out_channels = int(W*0.333), kernel_size=3, stride=1, padding=1, bias=False)
        self.bn5x5 = nn.BatchNorm2d(int(W*0.333))

        self.conv7x7 = nn.Conv2d(int(W*0.333), out_channels = int(W*0.5), kernel_size=3, stride=1, padding=1, bias=False)
        self.bn7x7 = nn.BatchNorm2d(int(W*0.5))

        self.shortcut = nn.Conv2d(in_channels, int(W*0.167) + int(W*0.333) + int(W*0.5), kernel_size=1, stride=1, padding=0, bias=False)
        self.bn_shortcut = nn.BatchNorm2d(int(W*0.167) + int(W*0.333) + int(W*0.5))

        self.bn_concat = nn.BatchNorm2d(int(W*0.167) + int(W*0.333) + int(W*0.5))
        self.bn_out = nn.BatchNorm2d(int(W*0.167) + int(W*0.333) + int(W*0.5))

    def forward(self, x):

        x_short = x
        x0 = self.bn_shortcut(self.shortcut(x_short))
        del x_short

        x1 =  F.relu(self.bn3x3(self.conv3x3(x)))
        x2 = F.relu(self.bn5x5(self.conv5x5(x1)))
        x3 = F.relu(self.bn7x7(self.conv7x7(x2)))
        x4 = self.bn_concat(torch.cat([x1, x2, x3], dim=1))
        del x1, x2, x3

        out = x0 + x4
        out = self.bn_out(F.relu(out))
        return out

class ResPath1(nn.Module):

    def __init__(self, in_channels, out_channels):
        super(ResPath1, self).__init__()

        self.shortcut0 = nn.Conv2d(in_channels, out_channels, kernel_size=1, stride=1, padding=0, bias=False)
        self.bn_shortcut0 = nn.BatchNorm2d(out_channels)

        self.shortcut1 = nn.Conv2d(out_channels, out_channels, kernel_size=1, stride=1, padding=0, bias=False)
        self.bn_shortcut1 = nn.BatchNorm2d(out_channels)

        self.shortcut2 = nn.Conv2d(out_channels, out_channels, kernel_size=1, stride=1, padding=0, bias=False)
        self.bn_shortcut2 = nn.BatchNorm2d(out_channels)

        self.shortcut3 = nn.Conv2d(out_channels, out_channels, kernel_size=1, stride=1, padding=0, bias=False)
        self.bn_shortcut3 = nn.BatchNorm2d(out_channels)

        self.conv0 = nn.Conv2d(in_channels, out_channels, kernel_size=3, stride=1, padding=1, bias=False)
        self.bn_conv0 = nn.BatchNorm2d(out_channels)

        self.conv1 = nn.Conv2d(out_channels, out_channels, kernel_size=3, stride=1, padding=1, bias=False)
        self.bn_conv1 = nn.BatchNorm2d(out_channels)

        self.conv2 = nn.Conv2d(out_channels, out_channels, kernel_size=3, stride=1, padding=1, bias=False)
        self.bn_conv2 = nn.BatchNorm2d(out_channels)

        self.conv3 = nn.Conv2d(out_channels, out_channels, kernel_size=3, stride=1, padding=1, bias=False)
        self.bn_conv3 = nn.BatchNorm2d(out_channels)

        self.bn_out0 = nn.BatchNorm2d(out_channels)
        self.bn_out1 = nn.BatchNorm2d(out_channels)
        self.bn_out2 = nn.BatchNorm2d(out_channels)
        self.bn_out3 = nn.BatchNorm2d(out_channels)

    def forward(self, x):

        x_short = x
        x_short0 = self.bn_shortcut0(self.shortcut0(x_short))
        x0 = self.bn_conv0(F.relu(self.conv0(x)))
        out = x_short0 + x0
        out = self.bn_out0(F.relu(out))

        x_short = out
        x_short1 = self.bn_shortcut1(self.shortcut1(x_short))
        x1 = self.bn_conv1(F.relu(self.conv1(out)))
        out = x_short1 + x1
        out = self.bn_out1(F.relu(out))

        x_short = out
        x_short2 = self.bn_shortcut2(self.shortcut2(x_short))
        x2 = self.bn_conv2(F.relu(self.conv1(out)))
        out = x_short2 + x2
        out = self.bn_out2(F.relu(out))

        x_short = out
        x_short3 = self.bn_shortcut3(self.shortcut3(x_short))
        x3 = self.bn_conv3(F.relu(self.conv3(out)))
        out = x_short3 + x3
        out = self.bn_out3(F.relu(out))

        return out

class ResPath2(nn.Module):

    def __init__(self, in_channels, out_channels):
        super(ResPath2, self).__init__()

        self.shortcut0 = nn.Conv2d(in_channels, out_channels, kernel_size=1, stride=1, padding=0, bias=False)
        self.bn_shortcut0 = nn.BatchNorm2d(out_channels)

        self.shortcut1 = nn.Conv2d(out_channels, out_channels, kernel_size=1, stride=1, padding=0, bias=False)
        self.bn_shortcut1 = nn.BatchNorm2d(out_channels)

        self.shortcut2 = nn.Conv2d(out_channels, out_channels, kernel_size=1, stride=1, padding=0, bias=False)
        self.bn_shortcut2 = nn.BatchNorm2d(out_channels)

        self.conv0 = nn.Conv2d(in_channels, out_channels, kernel_size=3, stride=1, padding=1, bias=False)
        self.bn_conv0 = nn.BatchNorm2d(out_channels)

        self.conv1 = nn.Conv2d(out_channels, out_channels, kernel_size=3, stride=1, padding=1, bias=False)
        self.bn_conv1 = nn.BatchNorm2d(out_channels)

        self.conv2 = nn.Conv2d(out_channels, out_channels, kernel_size=3, stride=1, padding=1, bias=False)
        self.bn_conv2 = nn.BatchNorm2d(out_channels)

        self.bn_out0 = nn.BatchNorm2d(out_channels)
        self.bn_out1 = nn.BatchNorm2d(out_channels)
        self.bn_out2 = nn.BatchNorm2d(out_channels)

    def forward(self, x):

        x_short = x
        x_short0 = self.bn_shortcut0(self.shortcut0(x_short))
        x0 = self.bn_conv0(F.relu(self.conv0(x)))
        out = x_short0 + x0
        out = self.bn_out0(F.relu(out))

        x_short = out
        x_short1 = self.bn_shortcut1(self.shortcut1(x_short))
        x1 = self.bn_conv1(F.relu(self.conv1(out)))
        out = x_short1 + x1
        out = self.bn_out1(F.relu(out))

        x_short = out
        x_short2 = self.bn_shortcut2(self.shortcut2(x_short))
        x2 = self.bn_conv2(F.relu(self.conv1(out)))
        out = x_short2 + x2
        out = self.bn_out2(F.relu(out))

        return out

class ResPath3(nn.Module):

    def __init__(self, in_channels, out_channels):
        super(ResPath3, self).__init__()

        self.shortcut0 = nn.Conv2d(in_channels, out_channels, kernel_size=1, stride=1, padding=0, bias=False)
        self.bn_shortcut0 = nn.BatchNorm2d(out_channels)

        self.shortcut1 = nn.Conv2d(out_channels, out_channels, kernel_size=1, stride=1, padding=0, bias=False)
        self.bn_shortcut1 = nn.BatchNorm2d(out_channels)

        self.conv0 = nn.Conv2d(in_channels, out_channels, kernel_size=3, stride=1, padding=1, bias=False)
        self.bn_conv0 = nn.BatchNorm2d(out_channels)

        self.conv1 = nn.Conv2d(out_channels, out_channels, kernel_size=3, stride=1, padding=1, bias=False)
        self.bn_conv1 = nn.BatchNorm2d(out_channels)

        self.bn_out0 = nn.BatchNorm2d(out_channels)
        self.bn_out1 = nn.BatchNorm2d(out_channels)

    def forward(self, x):

        x_short = x
        x_short0 = self.bn_shortcut0(self.shortcut0(x_short))
        x0 = self.bn_conv0(F.relu(self.conv0(x)))
        out = x_short0 + x0
        out = self.bn_out0(F.relu(out))

        x_short = out
        x_short1 = self.bn_shortcut1(self.shortcut1(x_short))
        x1 = self.bn_conv1(F.relu(self.conv1(out)))
        out = x_short1 + x1
        out = self.bn_out1(F.relu(out))

        return out


class ResPath4(nn.Module):

    def __init__(self, in_channels, out_channels):
        super(ResPath4, self).__init__()

        self.shortcut0 = nn.Conv2d(in_channels, out_channels, kernel_size=1, stride=1, padding=0, bias=False)
        self.bn_shortcut0 = nn.BatchNorm2d(out_channels)

        self.conv0 = nn.Conv2d(in_channels, out_channels, kernel_size=3, stride=1, padding=1, bias=False)
        self.bn_conv0 = nn.BatchNorm2d(out_channels)

        self.bn_out0 = nn.BatchNorm2d(out_channels)

    def forward(self, x):

        x_short = x
        x_short0 = self.bn_shortcut0(self.shortcut0(x_short))
        x0 = self.bn_conv0(F.relu(self.conv0(x)))
        out = x_short0 + x0
        out = self.bn_out0(F.relu(out))

        return out

class MultiResUNet(nn.Module):

    def __init__(self, n_channels, n_classes, patch_width, patch_height, alpha=1.67):
        super(MultiResUNet, self).__init__()

        self.n_channels = n_channels
        self.n_classes = n_classes
        self.patch_width = patch_width
        self.patch_height = patch_height

        self.mresblock1 = MultiResBlock(self.n_channels, 32)
        W = alpha * 32
        self.respath1 = ResPath1(int(W*0.167) + int(W*0.333) + int(W*0.5), 32)

        self.mresblock2 = MultiResBlock(int(W*0.167) + int(W*0.333) + int(W*0.5), 64)
        W = alpha * 64
        self.respath2 = ResPath2(int(W*0.167) + int(W*0.333) + int(W*0.5), 64)

        self.mresblock3 = MultiResBlock(int(W*0.167) + int(W*0.333) + int(W*0.5),128)
        W = alpha * 128
        self.respath3 = ResPath3(int(W*0.167) + int(W*0.333) + int(W*0.5),128)

        self.mresblock4 = MultiResBlock(int(W*0.167) + int(W*0.333) + int(W*0.5), 256)
        W = alpha * 256
        self.respath4 = ResPath4(int(W*0.167) + int(W*0.333) + int(W*0.5), 256)

        self.mresblock5 = MultiResBlock(int(W*0.167) + int(W*0.333) + int(W*0.5), 512)
        W = alpha * 512

        self.conv6_up = nn.ConvTranspose2d(int(W*0.167) + int(W*0.333) + int(W*0.5), 256, kernel_size=4, stride=2, padding=1)
        self.conv6_bn = nn.BatchNorm2d(256)
        self.mresblock6 = MultiResBlock(512, 256)
        W = alpha * 256

        self.conv7_up = nn.ConvTranspose2d(int(W*0.167) + int(W*0.333) + int(W*0.5), 128, kernel_size=4, stride=2, padding=1)
        self.conv7_bn = nn.BatchNorm2d(128)
        self.mresblock7 = MultiResBlock(256, 128)
        W = alpha * 128

        self.conv8_up = nn.ConvTranspose2d(int(W*0.167) + int(W*0.333) + int(W*0.5), 64, kernel_size=4, stride=2, padding=1)
        self.conv8_bn = nn.BatchNorm2d(64)
        self.mresblock8 = MultiResBlock(128, 64)
        W = alpha * 64

        self.conv9_up = nn.ConvTranspose2d(int(W*0.167) + int(W*0.333) + int(W*0.5), 32, kernel_size=4, stride=2, padding=1)
        self.conv9_bn = nn.BatchNorm2d(32)
        self.mresblock9 = MultiResBlock(64, 32)
        W = alpha * 32

        self.conv10 = nn.Conv2d(int(W*0.167) + int(W*0.333) + int(W*0.5), self.n_classes, kernel_size=3, stride=1, padding=1)

        self.pool = torch.nn.MaxPool2d(2)

    def forward(self, x):

        x1 = self.mresblock1(x)
        res_x1 = self.respath1(x1)
        x1 = self.pool(x1)

        x2 = self.mresblock2(x1)
        del x1
        res_x2 = self.respath2(x2)
        x2 = self.pool(x2)

        x3 = self.mresblock3(x2)
        del x2
        res_x3 = self.respath3(x3)
        x3 = self.pool(x3)

        x4 = self.mresblock4(x3)
        del x3
        res_x4 = self.respath4(x4)
        x4 = self.pool(x4)

        x5 = self.mresblock5(x4)
        del x4

        ### up
        up6 = F.relu(self.conv6_bn(self.conv6_up(x5)))
        up6 = torch.cat([up6, res_x4], dim=1)
        del x5, res_x4

        x6 = self.mresblock6(up6)
        del up6

        up7 = F.relu(self.conv7_bn(self.conv7_up(x6)))
        up7 = torch.cat([up7, res_x3], dim=1)
        del x6, res_x3

        x7 = self.mresblock7(up7)
        del up7

        up8 = F.relu(self.conv8_bn(self.conv8_up(x7)))
        up8 = torch.cat([up8, res_x2], dim=1)
        del x7, res_x2

        x8 = self.mresblock8(up8)
        del up8

        up9 = F.relu(self.conv9_bn(self.conv9_up(x8)))
        up9 = torch.cat([up9, res_x1], dim=1)
        del x8, res_x1

        x9 = self.mresblock9(up9)
        del up9

        x10 = self.conv10(x9)
        del x9
        x10 = x10.view(-1, self.n_classes, self.patch_width*self.patch_height)
        x10 = x10.permute((0, 2, 1))
        x10 = torch.sigmoid(x10)

        return x10
