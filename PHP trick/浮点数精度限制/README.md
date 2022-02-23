当小数小于`10^(-16)`后,PHP无法区分大小

```php
<?php
var_dump(1.000000000000001 > 1);
var_dump(1.0000000000000001 > 1);
```

```
boolean true
boolean false
```