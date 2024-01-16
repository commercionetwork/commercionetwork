
# CommercioAPI eSign - `Beta Version`

<!-- npm run docs:serve  -->

<!-- https://lcd-testnet.commercio.network/docs/did:com:1ug9j7hgaxu6mvfu2kgfdt3hqxn4mrwuztxc7nu/received -->


`In Review  - cooming soon`

## eSign


A Pades basic e-signature for a document is avaialable through API of the commercio app

https://dev-api.commercio.app/commercionetwork/v1/swagger/index.html#/Sign


The process is quite simple . You submit a file with a note to the Api 

You will get a :

* pdf Pades signed with the SDN (Subject distinguish name) in the metadata of the file of the user you are logged on and a selfsigned certificate
* xml file signed with a digital seal by commercio.network containing specific on pdf signature process

### PATH

POST : `/sign/process`


#### Step by step Example

##### Step 1 - Define the file to be signed  

For example you want to sign a pds document with your account `fw8ben.pdf`


#### Step 2 - Use the API to Send the message 


**Use the tryout**

COOMING SOON

**Corresponding Cli request**


```
curl -X 'POST' \
  'https://dev-api.commercio.app/v1/sign/process' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiI....u5Q' \
  -H 'Content-Type: multipart/form-data' \
  -F 'document=@fw8ben.pdf;type=application/pdf' \
  -F 'note=Order'
```

**API : Body response**

```

{
  "status": 1,
  "signature_type": 0,
  "retry_count": 0,
  "created_at": "2023-07-11T09:37:56.332391Z",
  "sign_start_at": "2023-07-11T09:37:56.402033Z",
  "sign_end_at": "0001-01-01T00:00:00Z",
  "tx_send_at": "0001-01-01T00:00:00Z",
  "id": "6ea92a56-184d-467c-94af-616bb6441ae8",
  "original_document_id": "7fa502ef-7565-4887-b207-c2f015cca5c8",
  "signed_document_id": "",
  "tx_process_id": "",
  "owner_id": "997d27ca-ee15-47a6-a29e-d239ca7050ac",
  "owner_username": "enterpriseuser001@zotsell.com",
  "owner_fullname": "Robert Plant",
  "owner_email": "enterpriseuser001@zotsell.com",
  "owner_dn": "C=IT,O=FoxSign,CN=Robert Plant",
  "signer_id": "997d27ca-ee15-47a6-a29e-d239ca7050ac",
  "signer_username": "enterpriseuser001@zotsell.com",
  "signer_fullname": "Robert Plant",
  "signer_email": "enterpriseuser001@zotsell.com",
  "signer_dn": "C=IT,O=FoxSign,CN=Robert Plant",
  "signer_certificate": "",
  "signer_certificate_chain": "",
  "note": "Accetto",
  "last_error": "",
  "original_document": {
    "id": "",
    "created_at": "",
    "size": 0,
    "label": "",
    "content_type": "",
    "original_name": "",
    "hash": "",
    "storage_uri": "",
    "description": ""
  },
  "signed_document": {
    "id": "",
    "created_at": "",
    "size": 0,
    "label": "",
    "content_type": "",
    "original_name": "",
    "hash": "",
    "storage_uri": "",
    "description": ""
  }
}
```

##### Step 3 - Check the process status 
YOur process id is "id": "6ea92a56-184d-467c-94af-616bb6441ae8",

Use the API Get by id:

**Use the tryout**



**Corresponding Cli request**


```
curl -X 'GET' \
  'https://dev-api.commercio.app/v1/sign/process/6ea92a56-184d-467c-94af-616bb6441ae8' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJSUz.....gu5Q'
```

**API : Body response**

```

{
  "status": 2,
  "signature_type": 0,
  "retry_count": 1,
  "created_at": "2023-07-11T09:37:56.332391Z",
  "sign_start_at": "2023-07-11T09:37:56.402033Z",
  "sign_end_at": "2023-07-11T09:38:02.776078Z",
  "tx_send_at": "2023-07-11T09:38:02Z",
  "id": "6ea92a56-184d-467c-94af-616bb6441ae8",
  "original_document_id": "7fa502ef-7565-4887-b207-c2f015cca5c8",
  "signed_document_id": "6397b211-2c51-4814-95d3-56768602b601",
  "tx_process_id": "70f26429-fae5-4948-a11e-9cbf73657ae5",
  "owner_id": "997d27ca-ee15-47a6-a29e-d239ca7050ac",
  "owner_username": "enterpriseuser001@zotsell.com",
  "owner_fullname": "Robert Plant",
  "owner_email": "enterpriseuser001@zotsell.com",
  "owner_dn": "C=IT,O=FoxSign,CN=Robert Plant",
  "signer_id": "997d27ca-ee15-47a6-a29e-d239ca7050ac",
  "signer_username": "enterpriseuser001@zotsell.com",
  "signer_fullname": "Robert Plant",
  "signer_email": "enterpriseuser001@zotsell.com",
  "signer_dn": "C=IT,O=FoxSign,CN=Robert Plant",
  "signer_certificate": "-----BEGIN CERTIFICATE-----\nMIIDbTCCAlWgAwIBAgIUGL0COgcwXzU.....y69gBrnAOBxmSCilOSYInNX/S540Rn1huNkyshw3HJNM8grePxDcHYepnjv+Z\nSYJX6PMShyuZZVJi49QFpFc=\n-----END CERTIFICATE-----\n",
  "signer_certificate_chain": "TUlBR0NTcUdTSWIzRFFFSEFxQ0FNSUFDQVFFeER6Q....UdOOENBQUFBQUFBQQ==",
  "note": "Order",
  "last_error": "",
  "original_document": {
    "id": "7fa502ef-7565-4887-b207-c2f015cca5c8",
    "created_at": "2023-07-11T09:37:56.332391Z",
    "size": 67700,
    "label": "",
    "content_type": "application/octet-stream",
    "original_name": "fw8ben.pdf",
    "hash": "1242833dff6c214973bd2bf902443133",
    "storage_uri": "documents/users/enterpriseuser001@zotsell.com/original/dc2c7d09-7148-4277-8413-312ef4465809",
    "description": ""
  },
  "signed_document": {
    "id": "6397b211-2c51-4814-95d3-56768602b601",
    "created_at": "2023-07-11T09:38:02.767873Z",
    "size": 567548,
    "label": "",
    "content_type": "application/octet-stream",
    "original_name": "signed-fw8ben.pdf",
    "hash": "10f00f46cf7122a9fb2c11b6a136f705",
    "storage_uri": "documents/users/enterpriseuser001@zotsell.com/signed/8847ef0b-9ef8-428b-808a-c18131ec878e",
    "description": ""
  }
}
```
##### Step 4 - Get your signed file 
YOur process id is "id": "6ea92a56-184d-467c-94af-616bb6441ae8",

Use the API Get by id: /sign/process/#id#/signed-document

**Use the tryout**



**Corresponding Cli request**

```
curl -X 'GET' \
  'https://dev-api.commercio.app/v1/sign/process/6ea92a56-184d-467c-94af-616bb6441ae8/signed-document' \
  -H 'accept: application/pdf' \
  -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJwSnpWTkVBa1JieGJvazJGajZPenlmR3RNR25IRVhYNjA4bEVD...HKbQ' \
  -o signed-fw8ben.pdf  
 ``` 

You wil get the file named signed-fw8ben.pdf

##### Step 5 - Get your audit file 


Use the API Get by id:

**Use the tryout**



**Corresponding Cli request**

```
  curl -X 'GET' \
  'https://dev-api.commercio.app/v1/sign/process/6ea92a56-184d-467c-94af-616bb6441ae8/audit' \
  -H 'accept: application/xml' \
  -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiI...ngfHKbQ' \
  -o audit.xml
```



You will obtain a Xml file xades signed with a digital seal issued by commercio.network containing basic data regarding the signature process of the pdf

```

  <?xml version="1.0" encoding="UTF-8"?>
<ds:Signature
	xmlns:ds="http://www.w3.org/2000/09/xmldsig#" Id="xmldsig-9d71af85-8720-4e11-95de-21287086bf01">
	<ds:SignedInfo>
		<ds:CanonicalizationMethod Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"/>
		<ds:SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"/>
		<ds:Reference Id="xmldsig-9d71af85-8720-4e11-95de-21287086bf01-ref0" Type="http://www.w3.org/2000/09/xmldsig#Object" URI="#xmldsig-9d71af85-8720-4e11-95de-21287086bf01-object0">
			<ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>
			<ds:DigestValue>COql3QGkOvrdJ/D74lRj21paiqE2yzcnQT9ZxCQjzQQ=</ds:DigestValue>
		</ds:Reference>
		<ds:Reference Type="http://uri.etsi.org/01903#SignedProperties" URI="#xmldsig-9d71af85-8720-4e11-95de-21287086bf01-signedprops">
			<ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>
			<ds:DigestValue>IwSgrClpmOZYEHCcrkAvS/xKuG86KdQ52yyKzZ7CXWg=</ds:DigestValue>
		</ds:Reference>
	</ds:SignedInfo>
	<ds:SignatureValue Id="xmldsig-9d71af85-8720-4e11-95de-21287086bf01-sigvalue">kuWfj8YLMMlUtEmZl65i4ArIXdV/aRvr0DwNKbaJkjscDAx08vVrgAxwD3aK73xX+JxEFKySYNlbHzvexpuO0xXQPJx7gyy1VrlXCe7+egY2oOvsPIzIid5yxf3tEqamyobpJ3KA+IxGEYDeeLmQeiVW41DfMZs/fxXKJ9PB/Fs7a5ih+gxDDI4Je7gZM1kxwQy9qzLt6ElhAS+ebyMfPVK2MUdR/Gym0/Af1JnvzqYTbAMrwQA/uZNVegSUlOK+nNRN/slW1S3q7u7hqFdbzMqwDm4weSrQ3KeyEHYPOcqTOuz7B2apGiVD8/qoSMqHZ2EI+RHQ27Uup5W9TpxwQw==</ds:SignatureValue>
	<ds:KeyInfo>
		<ds:X509Data>
			<ds:X509Certificate>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</ds:X509Certificate>
		</ds:X509Data>
	</ds:KeyInfo>
	<ds:Object Id="xmldsig-9d71af85-8720-4e11-95de-21287086bf01-object0">
		<Audit>
			<Document>
				<CreatedAt>2023-07-11T09:37:56.332391Z</CreatedAt>
				<ID>6ea92a56-184d-467c-94af-616bb6441ae8</ID>
				<Note>Accetto</Note>
				<OriginalDocument>
					<ContentType>application/octet-stream</ContentType>
					<CreatedAt>2023-07-11T09:37:56.332391Z</CreatedAt>
					<Hash>1242833dff6c214973bd2bf902443133</Hash>
					<ID>7fa502ef-7565-4887-b207-c2f015cca5c8</ID>
					<OriginalName>fw8ben.pdf</OriginalName>
					<Size>67700</Size>
					<StorageURI>documents/users/enterpriseuser001@zotsell.com/original/dc2c7d09-7148-4277-8413-312ef4465809</StorageURI>
				</OriginalDocument>
				<OriginalDocumentID>7fa502ef-7565-4887-b207-c2f015cca5c8</OriginalDocumentID>
				<OwnerDN>C=IT,O=FoxSign,CN=Robert Plant</OwnerDN>
				<OwnerEmail>enterpriseuser001@zotsell.com</OwnerEmail>
				<OwnerFullname>Robert Plant</OwnerFullname>
				<OwnerID>997d27ca-ee15-47a6-a29e-d239ca7050ac</OwnerID>
				<OwnerUsername>enterpriseuser001@zotsell.com</OwnerUsername>
				<RetryCount>1</RetryCount>
				<SignEndAt>2023-07-11T09:38:02.776078Z</SignEndAt>
				<SignStartAt>2023-07-11T09:37:56.402033Z</SignStartAt>
				<SignedDocument>
					<ContentType>application/octet-stream</ContentType>
					<CreatedAt>2023-07-11T09:38:02.767873Z</CreatedAt>
					<Hash>10f00f46cf7122a9fb2c11b6a136f705</Hash>
					<ID>6397b211-2c51-4814-95d3-56768602b601</ID>
					<OriginalName>signed-fw8ben.pdf</OriginalName>
					<Size>567548</Size>
					<StorageURI>documents/users/enterpriseuser001@zotsell.com/signed/8847ef0b-9ef8-428b-808a-c18131ec878e</StorageURI>
				</SignedDocument>
				<SignedDocumentID>6397b211-2c51-4814-95d3-56768602b601</SignedDocumentID>
				<SignerCertificate>-----BEGIN CERTIFICATE-----
MIIDbTCCAlWgAwIBAgIUGL0COgcwXzUA....hyuZZVJi49QFpFc=
-----END CERTIFICATE-----
</SignerCertificate>
				<SignerCertificateChain>TUlBR0NTcUdTS...QQ==</SignerCertificateChain>
				<SignerDN>C=IT,O=FoxSign,CN=Robert Plant</SignerDN>
				<SignerEmail>enterpriseuser001@zotsell.com</SignerEmail>
				<SignerFullname>Robert Plant</SignerFullname>
				<SignerID>997d27ca-ee15-47a6-a29e-d239ca7050ac</SignerID>
				<SignerUsername>enterpriseuser001@zotsell.com</SignerUsername>
				<Status>2</Status>
				<TxProcessID>70f26429-fae5-4948-a11e-9cbf73657ae5</TxProcessID>
				<TxSendAt>2023-07-11T09:38:02Z</TxSendAt>
			</Document>
		</Audit>
	</ds:Object>
	<ds:Object>
		<xades:QualifyingProperties
			xmlns:xades="http://uri.etsi.org/01903/v1.3.2#"
			xmlns:xades141="http://uri.etsi.org/01903/v1.4.1#" Target="#xmldsig-9d71af85-8720-4e11-95de-21287086bf01">
			<xades:SignedProperties Id="xmldsig-9d71af85-8720-4e11-95de-21287086bf01-signedprops">
				<xades:SignedSignatureProperties>
					<xades:SigningTime>2023-07-11T10:49:26.592Z</xades:SigningTime>
					<xades:SigningCertificate>
						<xades:Cert>
							<xades:CertDigest>
								<ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>
								<ds:DigestValue>M/EZgFQyXy/3vqfwswb9CQLQqAyNGpnPM3Z6b3nudBk=</ds:DigestValue>
							</xades:CertDigest>
							<xades:IssuerSerial>
								<ds:X509IssuerName>CN=GLOBALTRUST 2020 AATL 1,O=e-commerce monitoring GmbH,C=AT</ds:X509IssuerName>
								<ds:X509SerialNumber>66030374559097499702320121</ds:X509SerialNumber>
							</xades:IssuerSerial>
						</xades:Cert>
						<xades:Cert>
							<xades:CertDigest>
								<ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>
								<ds:DigestValue>62yvVlC0TiyWKq++AbQ1mxH3cyl4Ts1xWm7ENDfCtkw=</ds:DigestValue>
							</xades:CertDigest>
							<xades:IssuerSerial>
								<ds:X509IssuerName>CN=GLOBALTRUST 2020,O=e-commerce monitoring GmbH,C=AT</ds:X509IssuerName>
								<ds:X509SerialNumber>18514327146496444639718413</ds:X509SerialNumber>
							</xades:IssuerSerial>
						</xades:Cert>
						<xades:Cert>
							<xades:CertDigest>
								<ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>
								<ds:DigestValue>milqUYLR1FGi439Dm3Tar6JnUjMp+Q+aDSAHwzTiPJo=</ds:DigestValue>
							</xades:CertDigest>
							<xades:IssuerSerial>
								<ds:X509IssuerName>CN=GLOBALTRUST 2020,O=e-commerce monitoring GmbH,C=AT</ds:X509IssuerName>
								<ds:X509SerialNumber>109160994242082918454945253</ds:X509SerialNumber>
							</xades:IssuerSerial>
						</xades:Cert>
					</xades:SigningCertificate>
				</xades:SignedSignatureProperties>
			</xades:SignedProperties>
		</xades:QualifyingProperties>
	</ds:Object>
</ds:Signature>
```

### Enanchement
In order to make the signature procedure stronger You could also notarize one or both files in the chain (signed pdf and signed audit)

