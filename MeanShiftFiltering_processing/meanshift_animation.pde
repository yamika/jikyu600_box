int angleX;
int angleY;
int startX,startY;
float hs,hr;
PImage img;
PImage convert_img;
color[] plot_pixels;
color[] old_pyramid_pixels;
color[] output_pixels;

void setup(){
    size(800, 600, P3D);
    img = loadImage("lenna_rgb.jpg");
    img.resize(256, 256);
    convert_img = createImage(256, 256, RGB);
    plot_pixels = new color[img.width*img.height];
    old_pyramid_pixels = new color[img.width*img.height];
    output_pixels = new color[img.width*img.height];
    plot_pixels = img.pixels.clone();
    output_pixels = img.pixels.clone();
    angleX = 0; 
    angleY = 0;
    startX = -500;
    startY = -300;
    hs = 16;
    hr = 16;
}
 
void draw(){
    background(40);
    translate(500, 300);
    fill(255);
    textSize(16);
    text("hs : "+hs+"  hr : "+hr, 0, startY+50);
    image(img, startX, startY);
    image(convert_img, startX, startY+300);
    //rotateZ(radians(angleX));
    rotateY(radians(angleY));
    
    // 色空間用の立体描画
    stroke(255);
    noFill();
    strokeWeight(1);
    box(255, 255, 255);
    
    // 画像のピクセルを3Dプロット
    for (int i = 0; i < img.width*img.height; i++) {
        color c = plot_pixels[i];
        strokeWeight(2);
        stroke(c);
        point(red(c)-127,green(c)-127,blue(c)-127);
    }
}

class Pixel_vec {
  int x,y;
  float r,g,b;
    Pixel_vec(int x, int y,float r, float g, float b){
      this.x = x;
      this.y = y;
      this.r = r;
      this.g = g;
      this.b = b;
  }
   int get_x(){
      return x;
   }
   int get_y(){
      return y;
   }
   float get_r(){
      if(this.r < 0) this.r = 0.;
      if(this.r > 255.) this.r = 255.;
      return r;
   }
    float get_g(){
      if(this.g < 0) this.g = 0.;
      if(this.g > 255.) this.g = 255.;      
      return g;
   }
    float get_b(){
      if(this.b < 0) this.b = 0.;
      if(this.b > 255.)this.b = 255.;      
      return b;
   }
}

Pixel_vec get_center(PImage data, int x_dash, int y_dash, float r_dash, float g_dash, float b_dash){
    color c;
    float r = 0, g = 0, b = 0;
    
    float xy_distance = 0;
    float rgb_distance = 0;  
    
    float len = 0;  
    float sum_x = 0, sum_y = 0;
    float sum_r = 0, sum_g = 0, sum_b = 0;
    float center_x = 0, center_y = 0;
    float center_r = 0, center_g = 0, center_b = 0;
    
    Pixel_vec center_p;
    for(int x = 0; x < data.width; x++){
        for(int y = 0; y < data.height; y++){
          
            c = data.pixels[x*data.width+y];
            r = red(c);
            g = green(c);
            b = blue(c);
            
            //xy_distance = abs(x-x_dash) + abs(y-y_dash);
            rgb_distance = sqrt(pow(r-r_dash, 2) + pow(g-g_dash, 2) + pow(b-b_dash, 2));
            if(abs(x-x_dash) <= hs && abs(y-y_dash) <= hs && rgb_distance <= hr){
            //if(xy_distance <= hs && rgb_distance <= hr){
                len += 1;
                
                sum_x += x;
                sum_y += y;
                sum_r += r;
                sum_g += g;
                sum_b += b;
                
                center_x = sum_x / len;
                center_y = sum_y / len;
                center_r = sum_r / len;
                center_g = sum_g / len;
                center_b = sum_b / len;
            }
        }
  }
  center_p = new Pixel_vec(int(center_x), int(center_y), center_r, center_g, center_b);
  return center_p;
}

Boolean check_diff(color filter_c, color old_filter_c){
    float r,g,b;
    float old_r,old_g,old_b;
    
    r = red(filter_c);
    g = green(filter_c);
    b = blue(filter_c);
    old_r = red(old_filter_c);
    old_g = green(old_filter_c);
    old_b = blue(old_filter_c);
    
    if(sqrt(pow(r-old_r, 2) + pow(g-old_g, 2) + pow(b-old_b, 2)) < hr){
        return true;
    }
    return false;
}

void MeanShiftFiltering(){
    int max_level = 1;
    int x_dash = 0;
    int y_dash = 0;
    float r_dash, g_dash, b_dash = 0;
    float r = 0, g = 0, b = 0;
    float condition = 0;
    float eps = 1.0;
    color c;
    color output_color;
    Pixel_vec center_p = new Pixel_vec(0,0,255,100,100);
    PImage filter_img;
    
    for(int pyramid = max_level; pyramid >= 0; pyramid--){
        filter_img = img.get();
        filter_img.resize(img.width/int(pow(2,pyramid)), img.height/int(pow(2,pyramid)));
        filter_img.resize(img.width, img.height);
        for(int x = 0; x < img.width; x++){
            for(int y = 0; y < img.height; y++){
                x_dash = x;
                y_dash = y;
                // 前ステップの低解像度の画像と現ステップの画像の色の異なり具合を調べる
                //if(pyramid < max_level && check_diff(filter_img.pixels[x*img.width+y], old_pyramid_pixels[x*img.width+y])) break;
                //c = output_pixels[x*img.width+y];
                c = filter_img.pixels[x*img.width+y];
                r_dash = red(c);
                g_dash = green(c);
                b_dash = blue(c);
                // 以下のステップで中心の色が見つけられなかったら入力画像の色を代入する
                c = img.pixels[x*img.width+y];
                r = red(c);
                g = green(c);
                b = blue(c);         
                for(int iter=0; iter < 5; iter++){
                    center_p = get_center(filter_img, x_dash, y_dash, r_dash, g_dash, b_dash);
                    condition = abs(x_dash - center_p.x) + abs(y_dash - center_p.y) 
                    + pow(r_dash - center_p.r, 2) + pow(g_dash - center_p.g, 2) + pow(b_dash - center_p.b, 2);
                    if(condition <= eps || (int(x_dash) == int(center_p.x) && int(y_dash) == int(center_p.y))){
                      r = center_p.get_r();
                      g = center_p.get_g();
                      b = center_p.get_b();
                      break;
                    }
                    x_dash = center_p.get_x();
                    y_dash = center_p.get_y();
                    r_dash = center_p.get_r();
                    g_dash = center_p.get_g();
                    b_dash = center_p.get_b();
                }
                output_color = color(r, g, b);
                output_pixels[x*img.width+y] = output_color;
            }
        }
        old_pyramid_pixels = filter_img.pixels.clone();
    }
}

void keyPressed(){
    if(key == 'Z' || key == 'z'){
        println("Start MeanShiftFiltering");
        MeanShiftFiltering();
        for(int i=0; i<img.width*img.height; i++){
            convert_img.pixels[i] = output_pixels[i];
            plot_pixels[i] = output_pixels[i];
        }
        convert_img.updatePixels();
        println("Done");
    }
}

//色空間用の立体をグリグリ動かす
void mouseDragged(){
  
   if(mouseX < pmouseX){
       angleX += 1;
   }else{
       angleX -= 1;
   }
   
   if(mouseY < pmouseY){
       angleY += 1;
   }else{
       angleY -= 1;
   }
   
   if (angleX > 360 || angleY > 360) {
       angleX = 0;
       angleY = 0;
   }
}