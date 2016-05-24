# property_chaincode

[![Deploy to Bluemix](https://bluemix.net/deploy/button.png)](https://bluemix.net/deploy?repository=https://github.com/srini2win/property_chaincode.git)

Please find below detail for each operation

#Register:
	Curl command : curl -X POST -d @register.json http://211.29.157.13:3000/chaincode --header "Content-Type:application/json"
	Response: will add shortly

#Search all properties:
  a. CURL command : curl -X POST -d @search_all.json http:/ 211.29.157.13:3000/chaincode --header "Content-Type:application/json"
  b. Response:

#Search with specific property:
  a. CURL command : curl -X POST -d @search.json http://211.29.157.13:3000/chaincode --header "Content-Type:application/json"
  b. SearchTypes:  “ALL" OR "Folio_ID" OR "Address" OR "LegalOwner" OR "BeneficialOwner"
  c. Response:

