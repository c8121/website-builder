<!DOCTYPE html>
<html lang="de">
<head>
    <meta charset="utf-8"/>
</head>

<body>
    <?php echo "Test" ?>
    <?="Test"?>

    <div class="page-content">
        <div class="tpl Html Markdown Text Image">
            <h2>{{ Meta.Name }}</h2>
            <pre>Page: {{ Page.Name }}</pre>
            <pre>Page: {{ Meta.ContentType }}</pre>
        </div>
    </div>

</body>
</html>