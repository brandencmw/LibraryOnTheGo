<?php
    require './auth.php';
    
    if (isset($_COOKIE['Authorization']) && authenticated($_COOKIE['Authorization'])) {
    ?>
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Admin Panel</title>
        </head>
        <body>
        <div class="content">
                <h2>Authors</h2>
                <button id="authorPageButton">Go to authors page</button>

                <h2>Books</h2>
                <button id="bookPageButton">Go to books page</button>
        </div>
        <script src="redirects.js"></script>
        </body>
        </html>
    <?php
    } else {
        header('Location: ../login');
    }
?>