<?php
function getBlock(){
$str = <<<EOF
- type: log
  enabled: true
  paths:
    - #paths#

  json.keys_under_root: true
  json.add_error_key: true

  fields:
    source: #source#
    
EOF;
    return $str;
}

function getIndex(){
$str = <<<EOF
    - index: "ck-local-#index#-%{+yyyy.MM.dd}"
      when.equals:
        fields:
            source: "#index#"
            
EOF;
    return $str;
}