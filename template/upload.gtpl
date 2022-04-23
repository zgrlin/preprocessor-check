<!DOCTYPE html>
<html lang="tr">
  <head>

<html lang="tr">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>Goosey - Ana Sayfa</title>

   <style>
        h3 {
            text-align: center;
        }
  
        img {
            display: block;
            margin: auto;
            height: 150px;
            width: 150px;
        }
  
        .input {
            margin: 6px;
            padding: 10px;
            display: block;
            margin: auto;
            color: palevioletred;
            font-size: 30px;
        }
  
        input {
            width: 90%;
            display: block;
            margin-left: 12px;
            background: none;
            background-color: lightyellow;
        }
  
        select {
            width: 90%;
            display: block;
            margin-left: 12px;
            background: none;
            background-color: lightyellow;
        }
  
        #heading {
            font-family: cursive;
            text-align: center;
            color: green;
            padding-top: 20px;
  
        }
  
        #form_page {
            height: 500px;
            width: 50%;
            display: flex;
            flex-wrap: wrap;
            flex-direction: row;
            margin: auto;
  
        }
  
        #form_body {
            border-radius: 12px;
            height: 330px;
            width: 450px;
            background-color: beige;
            border: 1px solid pink;
            margin: auto;
            margin-top: 12px;
        }
  
        #text {
            color: red;
            width: 100px;
        }
  
        #head {
            border-bottom: 2px solid red;
            height: 100px;
            background-color: aliceblue;
        }
  
        #submit {
            background-color: white;
            width: 70px;
        }
    </style>
  </head>
  <body>
     <button onclick="location.href='https://url/process'" type="button">Kontrol et</button>
     <button onclick="location.href='https://url/stop'" type="button">Durdur</button>
     <button onclick="location.href='https://url'" type="button">Calistir</button>
     <button onclick="location.href='https://url/log/'" type="button">Log dosyasini indir</button>
     <button onclick="location.href='https://url/logout'" type="button" style="float: right;">Cikis yap</button>
    <br>
    <h3> Dosya Yukleme </h3>
    <form id="form" enctype="multipart/form-data" action="https://url/upload" method="POST">
      <input class="input file-input" type="file" name="file" multiple />
      <button class="button" type="submit" style="position: absolute; margin: auto; left: 50%; top:50%;">Gonder gitsin</button>
    </form>
  </body>
</html>
