package com.example.yamikachan.roadsign_detection;

import android.Manifest;
import android.content.res.AssetManager;
import android.content.pm.PackageManager;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.graphics.Color;
import android.os.Bundle;
import android.os.Build;
import android.os.Environment;
import android.support.v4.app.ActivityCompat;
import android.support.v7.app.AppCompatActivity;
import android.util.Log;

import android.widget.ImageView;
import android.widget.TextView;
import android.widget.Toast;
import android.view.SurfaceView;
import android.view.WindowManager;
import org.tensorflow.contrib.android.TensorFlowInferenceInterface;

import org.opencv.android.OpenCVLoader;
import org.opencv.android.BaseLoaderCallback;
import org.opencv.android.CameraBridgeViewBase.CvCameraViewFrame;
import org.opencv.android.LoaderCallbackInterface;
import org.opencv.core.Mat;
import org.opencv.core.Core.MinMaxLocResult;
import org.opencv.core.CvType;
import org.opencv.core.Core;
import org.opencv.core.Rect;
import org.opencv.core.Scalar;
import org.opencv.core.Point;
import org.opencv.imgproc.Imgproc;
import org.opencv.android.CameraBridgeViewBase;
import org.opencv.android.CameraBridgeViewBase.CvCameraViewListener2;
import org.opencv.android.Utils;

public class MainActivity extends AppCompatActivity implements CvCameraViewListener2{
    static {
        System.loadLibrary("opencv_java3");
    }

    private static final String TAG = "OCVSample::Activity";
    private CameraBridgeViewBase mOpenCvCameraView;

    private final int REQUEST_PERMISSION = 1000;

    private TextView textView;
    private Bitmap bitmap_stop;
    private AssetManager mAssetManager;
    private TensorFlowInferenceInterface inferenceInterface;

    private BaseLoaderCallback mLoaderCallback = new BaseLoaderCallback(this) {
        @Override
        public void onManagerConnected(int status) {
            switch (status) {
                case LoaderCallbackInterface.SUCCESS:
                {
                    Log.i(TAG, "OpenCV loaded successfully");
                    mOpenCvCameraView.enableView();
                } break;
                default:
                {
                    super.onManagerConnected(status);
                } break;
            }
        }
    };

    public MainActivity() {
        Log.i(TAG, "Instantiated new " + this.getClass());
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {

        Log.i(TAG, "called onCreate");
        super.onCreate(savedInstanceState);
        getWindow().addFlags(WindowManager.LayoutParams.FLAG_KEEP_SCREEN_ON);

        setContentView(R.layout.activity_main);
        mOpenCvCameraView = (CameraBridgeViewBase) findViewById(R.id.camera_view);
        textView = (TextView) findViewById(R.id.textview);
        mOpenCvCameraView.setVisibility(SurfaceView.VISIBLE);
        if (Build.VERSION.SDK_INT >= 23) {
            // protocol buffer ファイルの読み込み
            checkPermission();
        }
        mOpenCvCameraView.setCvCameraViewListener(this);


    }

    private int Predict(Bitmap bitmap){

        // CNNにかける入力の前処理
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
                //RGB各色に分離して入力データの配列に格納
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
        // 入力とkeep_probの設定
        inferenceInterface.feed(input_x,input_data,1,224*224*3);
        inferenceInterface.feed(keep_prob,prob,1,1);
        String output = "output2:0";
        inferenceInterface.run(new String[]{output});
        final float[] outputs = new float[4];
        inferenceInterface.fetch(output,outputs);

        // 要素が最大のインデックスを返す
        return argmax(outputs);
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

            mAssetManager = getAssets();
            inferenceInterface = new TensorFlowInferenceInterface(mAssetManager, filePath);

        }
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

    @Override
    public void onPause()
    {
        super.onPause();
        if (mOpenCvCameraView != null)
            mOpenCvCameraView.disableView();
    }

    @Override
    public void onResume()
    {
        super.onResume();
        if (!OpenCVLoader.initDebug()) {
            Log.d(TAG, "Internal OpenCV library not found. Using OpenCV Manager for initialization");
            OpenCVLoader.initAsync(OpenCVLoader.OPENCV_VERSION_3_2_0, this, mLoaderCallback);
        } else {
            Log.d(TAG, "OpenCV library found inside package. Using it!");
            mLoaderCallback.onManagerConnected(LoaderCallbackInterface.SUCCESS);
        }
    }

    public void onDestroy() {
        super.onDestroy();
        if (mOpenCvCameraView != null)
            mOpenCvCameraView.disableView();
    }

    public void onCameraViewStarted(int width, int height) {
    }

    public void onCameraViewStopped() {
    }

    public Mat onCameraFrame(CvCameraViewFrame inputFrame) {
        // カメラから得られるフレームを取得
        Mat mRgba = inputFrame.rgba();
        Mat OutImg = new Mat();
        OutImg = TemplateMatching(mRgba,4);
        try {
            //Bitmapに変換する前にRGB形式に変換
            Imgproc.cvtColor(OutImg, OutImg, Imgproc.COLOR_RGBA2RGB, 4);
            //空のBitmapを作成
            Bitmap dst = Bitmap.createBitmap(OutImg.width(), OutImg.height(), Bitmap.Config.ARGB_8888);
            //MatからBitmapに変換
            Utils.matToBitmap(OutImg, dst);
            //CNNの入力サイズである 224x224 サイズに変換
            dst = Bitmap.createScaledBitmap(dst, 224, 224, false);
            //予測を行い、インデックスを取得
            final int idx = Predict(dst);
            runOnUiThread(new Runnable() {
                @Override
                public void run() {
                    //取得したインデックスから正解ラベルを表示
                    textView.setText(convert_idx_into_Label(idx));
                }
            });
        } catch (Exception e) {
            // TODO: handle exception
            e.printStackTrace();
        }
        return mRgba;
    }

    public Mat TemplateMatching(Mat targetImg,int match_method) {
        // テンプレート画像
        bitmap_stop = BitmapFactory.decodeResource(getResources(), R.drawable.stop);
        Mat templImg = new Mat();
        // bitmapからMax形式に変換
        Utils.bitmapToMat(bitmap_stop,templImg);

        // / Create the result matrix
        int result_cols = targetImg.cols() - templImg.cols() + 1;
        int result_rows = targetImg.rows() - templImg.rows() + 1;
        Mat result = new Mat(result_rows, result_cols, CvType.CV_32FC1);

        // / Do the Matching and Normalize
        Imgproc.matchTemplate(targetImg, templImg, result, match_method);
        Core.normalize(result, result, 0, 1, Core.NORM_MINMAX, -1, new Mat());

        // / Localizing the best match with minMaxLoc
        MinMaxLocResult mmr = Core.minMaxLoc(result);

        Point matchLoc;
        if (match_method == Imgproc.TM_SQDIFF || match_method == Imgproc.TM_SQDIFF_NORMED) {
            matchLoc = mmr.minLoc;
        } else {
            matchLoc = mmr.maxLoc;
        }
        // マッチングしたところを四角で囲む
        Imgproc.rectangle(targetImg, matchLoc, new Point(matchLoc.x + templImg.cols(),
                matchLoc.y + templImg.rows()), new Scalar(0, 255, 0),5);
        // マッチングしたところを切り抜く(CNNに渡す入力用)
        Rect roi = new Rect(matchLoc,new Point(matchLoc.x + templImg.cols(), matchLoc.y + templImg.rows()));
        Mat trim = new Mat(targetImg,roi);
        return trim;
    }
}

