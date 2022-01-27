# Amazon-S3-Bucket-Operations

~~~go
if len(os.Args) != 2 {
		exitErrorf("bucket name required\nUsage: %s bucket_name", os.Args[0])
	}
~~~

os.Args komut satırına girilen argümanları hafızada tutmaya yarar eğer gönderilen argüman 2 ye eşit değil ise hata versin 


~~~go
sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

// Create S3 service client	
svc := s3.New(sess)
~~~

ilk olarak default config ile bir oturum oluşturuyoruz burdaki veriler değiştirlmediği sürece aws nin kendi default deneme oturumu oluşmuş olur (Code aws doc tan alınmıştır)

> oluşan svc(*s3) istemcisi(nesnesi) ile yapılmak istenen **Buckets** işlemleri sunulur 

veya oluşan oturum ile oluşturulacak işlem objeleri olabilir 
ne demek istedik hemen bakalım 

~~~go 
	uploader := s3manager.NewUploader(sess)

    downloader := s3manager.NewDownloader(sess)
~~~

method isimlerinden anlaşılacagı üzere de ne yapılacagı hakkında bilgi verior 

sonrasında yapılacak işlemlerin kontrolleri ve işlem çıktıları mevcut 

[Daha detaylı bilgi için](!https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-basic-bucket-operations.html)