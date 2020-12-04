## Product Discount From xlsx 

This is a small quick utility built for converting a client xlsx file into a magento readable discount csv 

If you have go installed on your computer you can simply add your file named 
```discount_item_list.xlsx``` to a ```/data``` directory at the root of this project and use 

```go run main.go```

it will produce a translated csv file ready for magento 2 customer group discount import 

You can also build the executable of this project using ```go build``` or use the included linux executable if you are running on an amd64 linux based os 


### structure for the xsls file 
| sku | wholesaler | retailer |
|-|-|-|
| 1099 | 30 | 10 |

sku is product sku each column represents a given customer group, and the value in a non-sku cell represents the percentage discount that group receieves