<?php declare(strict_types=1);
require './get_all_authors.php';

function formatLongTextForCard(string $text, int $len) {
    if (strlen($text) <= $len) {
        return $text;
    }

    $text = substr($text,0, $len);
    return trim($text) . '...';
}

$order = isset($_GET['order']) ? $_GET['order'] : '';
$name = isset($_GET['name']) ? $_GET['name'] : '';
$authors = getAuthors($order, $name);

$sortOptions = array(
    'first_name_asc' => 'First Name A-Z',
    'first_name_desc' => 'First Name Z-A',
    'last_name_asc' => 'Last Name A-Z',
    'last_name_desc' => 'Last Name Z-A',
);
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
            <div class="filter-accordion" id="filterAccordion">
                <div class="filter-accordion-header" id="filterAccordionHeader">
                    <h3>Sort & Filter Results</h3>
                    <svg width="16" height="12" viewBox="0 0 16 12" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path d="M8 12L0.205772 0.299999L15.7942 0.299999L8 12Z" fill="#C4C4C4"/>
                    </svg>
                </div>

                <div class="filter-accordion-content" id="filterAccordionContent">
                    <select name="" id="sortDropdown">
                        <option value="default">Sort Results</option>
                        <?php foreach ($sortOptions as $value => $title) { ?>
                            <option value="<?php echo $value ?>" <?php if ($value == $order) {echo "selected";} ?>><?php echo $title ?></option>
                        <?php } ?>
                    </select>
                    <div class="filter-box">
                        <h4>Filter By Last Initial</h4>
                        <div class="alphabet-buttons">
                        <?php
                        foreach(range('A','Z') as $letter) {
                        ?>
                            <button class="alphabet-button<?php if ($letter == $name) {echo ' selected';} ?>" id="filter<?php echo $letter ?>"><?php echo $letter ?></button>
                        <?php
                        }
                        ?>
                        </div>
                    </div>
                </div>
            </div>
            <div class="browse-content-box">
                <?php if (sizeof($authors) == 0) { ?>
                    <h3>We were not able to find any authors with the given criteria</h3>
                <?php  } else {
                    foreach ($authors as $author) { ?>
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
                <?php
                    } 
                } ?>
            </div>
        </div>
    </div>
</body>
<script src="filter-accordion.js"></script>
</html>