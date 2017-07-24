# ROADSIGN Classification
道路標識の「制限速度」と「止まれ」の判別をする学習モデル    

### ROADSIGN Classification
全結合層のみのネットワーク

### ROADSIGN Classification ConvNet
畳み込み層を追加したネットワーク

### ROADSIGN Classification AlexNet
CNNの一つであるAlexNetを実装

### ROADSIGN Classification RCNN
SelectiveSearchとConvNetの方のネットワークを使って物体検出          
こちらはNon-Maximum-Suppressionを適用せず、一番得点の高いものを一つ検出するだけ       
SelectiveSearchのアルゴリズムは[Alpaca](https://github.com/AlpacaDB/selectivesearch)さんのライブラリを利用させていただきました。      

### template matching     
テンプレートマッチングによる物体の領域をAlexNetの学習モデルにかけ物体の判定を行う物体検出もどき     

### RCNN-Non-Maximum-Suppression      
SelectiveSearchとAlexNetの学習モデル、Non-Maximum-Suppressionを実装した物体検出      
