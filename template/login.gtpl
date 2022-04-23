<!DOCTYPE html>
 <html lang="tr">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>Goosey - Giris</title>

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
	<h3>Merhaba</h3>
        <form method="post" action="https://url/login">
	   <label for="username">Kullanici Adi: </label>
           <input type="text" id="username" name="username" placeholder="kullanici adiniz" required>
           <label for="password">Sifre: </label>
           <input type="password" id="password" name="password" placeholder="sifre lutfen" required>
           <button type="submit" name="submit" value="submit">Sansini dene!</button>
        </form>
    </body>
</html>
