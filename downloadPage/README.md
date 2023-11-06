sudo apt-get update && \
    sudo apt-get -y install locales fonts-ipafont fonts-ipaexfont && \
    echo "ja_JP UTF-8" > /etc/locale.gen && locale-gen


日本語文字化け対策
    https://qiita.com/skokado/items/a2a422c8636a919ce52f
