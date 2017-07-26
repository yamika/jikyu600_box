package com.example.yamikachan.roadsign_classification;

import android.Manifest;
import android.content.res.AssetManager;
import android.content.pm.PackageManager;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.graphics.Color;
import android.os.Bundle;
import android.os.Environment;
import android.os.Build;
import android.provider.Settings;
import android.support.v4.app.ActivityCompat;
import android.support.v7.app.AppCompatActivity;
import android.util.Log;

import android.widget.Button;
import android.widget.ImageView;
import android.widget.TextView;
import android.widget.Toast;
import android.view.View;


import org.tensorflow.contrib.android.TensorFlowInferenceInterface;
public class MainActivity extends AppCompatActivity {

    private final int REQUEST_PERMISSION = 1000;

    private TextView textView;
    private ImageView imageView;
    private Bitmap bitmap;
    private Button predict_button;

    private AssetManager mAssetManager;
    private TensorFlowInferenceInterface inferenceInterface;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        textView = (TextView) findViewById(R.id.textview);
        imageView = (ImageView) findViewById(R.id.imageview);
        predict_button = (Button) findViewById(R.id.predict_button);
        predict_button.setEnabled(false);

        bitmap = BitmapFactory.decodeResource(getResources(), R.drawable.stop);
        bitmap = Bitmap.createScaledBitmap(bitmap, 224, 224, false);
        imageView.setImageBitmap(bitmap);
        if (Build.VERSION.SDK_INT >= 23) {
            checkPermission();
        }

        predict_button.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
            Predict(bitmap);
            }
        });

    }

    private void Predict(Bitmap bitmap){
        int width = bitmap.getWidth();
        int height = bitmap.getHeight();
        int i = 0;
        int pixels[] = new int[width * height];
        float input_data[] = new float[width*height*3];
        float prob[] = new float[1];
        prob[0] = 1;
        bitmap.getPixels(pixels, 0, width, 0, 0, width, height);
        for (int y = 0; y < height; y++) {
            for (int x = 0; x < width; x++) {
                int pixel = pixels[x + y * width];
                //RGB各色に分離
                int r = Color.red(pixel);
                int g = Color.green(pixel);
                int b = Color.blue(pixel);
                input_data[i+0] = r / 255.f;
                input_data[i+1] = g / 255.f;
                input_data[i+2] = b / 255.f;
                i += 3;
            }
        }

        String input_x = "input:0";
        String keep_prob = "keep_prob:0";
        inferenceInterface.feed(input_x,input_data,1,224*224*3);
        inferenceInterface.feed(keep_prob,prob,1,1);
        String output = "output2:0";
        inferenceInterface.run(new String[]{output});
        final float[] outputs = new float[4];
        inferenceInterface.fetch(output,outputs);
        String predict_label = convert_idx_into_Label(argmax(outputs));
        textView.setText("This image is "+predict_label);
    }

    // 配列の要素が最大のインデックスを取得する
    private int argmax(float[] array) {
        int argmax=0;
        float maxvalue=array[0];
        System.out.println(maxvalue);
        for (int i=1; i<array.length; i++) {
            System.out.println(array[i]);
            if (maxvalue < array[i]) {
                argmax=i;
                maxvalue = array[i];
            }
        }
        return argmax;
    }

    // インデックスを正解ラベルに変換
    private String convert_idx_into_Label(int idx) {
        switch(idx){
            case 0:
                return "Stop";
            case 1:
                return "LimitSpeed10";
            case 2:
                return "LimitSpeed20";
            case 3:
                return "LimitSpeed30";
        }
        return "Failed";
    }

    // permissionの確認
    public void checkPermission() {
        // 既に許可している
        if (ActivityCompat.checkSelfPermission(this,
                Manifest.permission.WRITE_EXTERNAL_STORAGE)==PackageManager.PERMISSION_GRANTED){
            setUpReadExternalStorage();
        }
        // 拒否していた場合
        else{
            requestLocationPermission();
        }
    }
    private void setUpReadExternalStorage(){

        // 現在ストレージが読出しできるかチェック
        if(isExternalStorageReadable()){

            // ProtocolBufferのパスを得る
            String filePath = Environment.getExternalStorageDirectory().getPath()
                    + "/" + "frozen_model.pb";
            Log.d("debug", "filePath="+filePath);

            try {
                mAssetManager = getAssets();
                inferenceInterface = new TensorFlowInferenceInterface(mAssetManager, filePath);
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
        predict_button.setEnabled(true);
    }

    public boolean isExternalStorageReadable() {
        String state = Environment.getExternalStorageState();
        if (Environment.MEDIA_MOUNTED.equals(state) ||
                Environment.MEDIA_MOUNTED_READ_ONLY.equals(state)) {
            return true;
        }
        return false;
    }

    // 許可を求める
    private void requestLocationPermission() {
        if (ActivityCompat.shouldShowRequestPermissionRationale(this,
                Manifest.permission.WRITE_EXTERNAL_STORAGE)) {
            ActivityCompat.requestPermissions(MainActivity.this,
                    new String[]{Manifest.permission.WRITE_EXTERNAL_STORAGE}, REQUEST_PERMISSION);

        } else {
            Toast toast = Toast.makeText(this, "アプリ実行に許可が必要です", Toast.LENGTH_SHORT);
            toast.show();

            ActivityCompat.requestPermissions(this,
                    new String[]{Manifest.permission.WRITE_EXTERNAL_STORAGE,}, REQUEST_PERMISSION);

        }
    }
}
