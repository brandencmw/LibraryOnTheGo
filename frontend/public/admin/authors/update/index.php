<?php
require '../../auth.php';
function getAuthorToUpdate($id) {
    $update_endpoint = 'https://server/authors?id=' . $id;
    $ch = curl_init($update_endpoint);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_CUSTOMREQUEST,'GET');
    curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, true);
    curl_setopt($ch, CURLOPT_CAINFO, '/etc/ssl/certs/root-ca.pem');
    $json_response = curl_exec($ch);
    if (curl_errno($ch)) {
        echo ''. curl_error($ch);
        return null;
    }
    curl_close($ch);
    $response = json_decode($json_response, true);
    return $response['author'];

}

if (isset($_GET['id']) && isset($_COOKIE['Authorization']) && authenticated($_COOKIE['Authorization'])) {
    $author = getAuthorToUpdate($_GET['id']);
    if ($author != null) {
?>

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Update Author Record</title>
    <script src="../../send-data.js"></script>
</head>
<body>
    <h2>Update Author</h2>
    <form id="updateAuthorForm">
        <label for="authorFirstName">First Name</label>
        <input type="text" name="authorFirstName" id="authorFirstName" value="<?php echo $author['firstName']; ?>">

        <label for="authorLastName">Last Name</label>
        <input type="text" name="authorLastName" id="authorLastName" value="<?php echo $author['lastName']; ?>">

        <label for="authorBio">Bio</label>
        <textarea name="authorBio" id="authorBio"><?php echo $author['bio']; ?></textarea>

        <label for="authorHeadshot">Headshot</label>
        <input type="file" name="authorHeadshot" id="authorHeadshot">
        <img src="https://library-pictures.s3.amazonaws.com/<?php echo $author['headshotKey']; ?>" id="headshotPreview">
        <button type="submit">Submit</button>
    </form>
    <button id="backButton">Back</button>
</body>
<script src="update-author.js"></script>
</html>

<?php
    } else {
        header('Location: ../');
    }
} else { 
    header('Location: ../../../../login/');
}
?>