let TAX_RATE = 0.19
define net_price(gross_price) = gross_price / (1 + TAX_RATE)
define discount(price, percent) = price * (1 - percent / 100)
define final_price(price, percent) = net_price(discount(price, percent))