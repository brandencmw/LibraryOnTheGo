<?php declare(strict_types=1);
require './get_all_authors.php';

function formatLongTextForCard(string $text, int $len) {
    if (strlen($text) <= $len) {
        return $text;
    }

    $text = substr($text,0, $len);
    return trim($text) . '...';
}

$authors = getAuthors();
?>

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="stylesheet" href="../styles.css">
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Merriweather&family=Poppins&display=swap" rel="stylesheet">
    <title>System Authors</title>
</head>
<body>
    <div class="main-box">
        <?php include '../navbar.html' ?>
        <div class="results-frame">
            <div class="browse-content-box">
                <?php foreach ($authors as $author) { ?>
                    <div class="browse-card">
                        <div class="browse-card-img">
                            <img class="result-card-headshot" src="<?php echo 'https://library-pictures.s3.amazonaws.com/' . $author['headshotKey']; ?>" alt="">
                        </div>
                        <div class="browse-card-content">
                            <h2 class="result-card-title"><?php echo $author['firstName'] . ' ' . $author['lastName']; ?></h2>
                            <p class="result-card-body"><?php echo formatLongTextForCard($author['bio'], 150); ?></p>
                            <a href="/authors?id=<?php echo $author['id']; ?>">See More</a>
                        </div>
                    </div>
                <?php } ?>
            </div>
        </div>
    </div>
</body>
</html>