<?php
/**
 * Created by PhpStorm.
 * User: admor
 * Date: 06.10.14
 * Time: 14:44
 */

spl_autoload_register(
    function($className) {
        require_once str_replace('\\', '/', $className) . '.php';
    }
);

print 123;