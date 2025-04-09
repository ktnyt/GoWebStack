// ブラウザリロード用のWebSocket接続を確立する
(function() {
  console.log('リロード監視を開始します...');
  
  // WebSocket接続を作成
  const socket = new WebSocket('ws://localhost:6641/subscribe');
  
  // 接続が開いたときの処理
  socket.onopen = function() {
    console.log('リロード監視サーバーに接続しました');
  };
  
  // メッセージを受信したときの処理
  socket.onmessage = function(event) {
    if (event.data === 'reload') {
      console.log('リロード信号を受信しました。ページをリロードします...');
      window.location.reload();
    }
  };
  
  // エラーが発生したときの処理
  socket.onerror = function(error) {
    console.error('WebSocket接続エラー:', error);
  };
  
  // 接続が閉じたときの処理
  socket.onclose = function() {
    console.log('リロード監視サーバーとの接続が閉じられました。再接続を試みます...');
    // 少し待ってから再接続を試みる
    setTimeout(function() {
      window.location.reload();
    }, 2000);
  };
})();
