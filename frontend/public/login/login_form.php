<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login</title>
</head>
<body>
    <h2>Login to Admin Panel</h2>
    <form action="login.php" method="post" id="loginForm">
        <label for="username">Enter Username</label>
        <input type="text" name="username">

        <label for="password">Enter Password</label>
        <input type="password" name="password">
        <button type="submit">Submit</button>
    </form>
    <div id="errorBox"></div>
</body>
</html>