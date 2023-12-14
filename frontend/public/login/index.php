<?php
    require '../admin/auth.php';

    if (isset($_COOKIE['Authorization']) && authenticated($_COOKIE['Authorization'])) {
        header('Location: ../admin');
        exit();
    }
    include 'login_form.php';
?>