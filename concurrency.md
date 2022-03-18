  ## GO'DA CONCURRENCY            
             
Kodların tek bir rutin içerisinde yukarıdan aşağı satır satır okunarak çalışması çoğu zaman verimlilikten uzak bir sonuç verir. Bu duruma çözüm getirmek için kullanılan benzer fakat temelde farklı 2 kavram var: Concurrency ve Parallelism. Parallelism aynı anda birden fazla işlemin farklı çekirdeklerde gerçekleştirilmesi durumu iken, concurrency birbirinden bağımsız işlemlerin sonuçlarının etkilenmeden ve sıra gözetmeksizin gerçekleştirilmesidir. Aşağıdaki görselin konunun algılanması açısından oldukça faydalı olduğunu düşünüyorum.            

![concurrency vs paralellism](https://golangbot.com/content/images/2017/06/concurrency-parallelism-copy.png)          
             
           
Go Concurrency’de thread yerini goroutine almaktadır. Goroutine, thread’e göre çok daha hafifir. Goroutine’ler tıpkı fonksiyonlar gibi tanımlanır ve her biri Go Schedule içerisindeki bir rutini temsil eder. Birbirinden bağımsız işler yapan, fakat birbiriyle iletişim kurabilen iş parçacıkları olarak düşünebiliriz. Main fonksiyonu da bir goroutine oluşturur.
Basit bir örnekle başlayalım.      

![concurrency vs paralellism](concurrency.png)          
                   
printSomeText ve printOtherTexts fonksiyonları belli süreler bekleyerek bazı metinleri ekrana yazdırıyorlar. Concurrency kullanmadan bu uygulama çalıştırıldığında elbette önce ekranda sırayla metinlerin yazdırıldığını görürüz. Ne var ki programdan isteğimiz bu olmayabilir. Sıra fark etmeksizin zaman anlamında en hızlı sonucun verilmesini isteyebiliriz. İşte bu noktada concurrency devreye girmekte. Yukarıdaki uygulama çalıştığında beklentimiz metinlerin sıra beklentisi olmadan yazdırılması ve 1 saniyelik beklemelerin görmezden gelinmesidir elbette. Ancak bu örnekte sonuç olarak hiçbir şeyin yazdırılmadığını görürüz. Bunun sebebi ise main fonksiyonun da bağımsız bir go routine olduğunu unutmamızdır. Şu anlık normalde pek de tercih edilmemesi gereken bir çözüm olarak uygulamanın bir süre beklemesi için main fonksiyonuna sleep metodunu ekleyerek tekrar çalıştırırsak:     
      
![concurrency vs paralellism](concurrency2.png)                
      
![concurrency vs paralellism](concurrency3.png)                 
        
 
Çıktıda görüldüğü üzere printOtherTexts fonksiyonunu daha sonra çağırmamıza rağmen, önce o rutinde bulunan metinleri yazdı. Geçen zaman da main fonksiyonu sonunda yarım saniyelik bir bekleme olmasını istediğimiz için 500 ms civarında.
Uygulamanın durumuna göre text 3’den hemen sonra text 4’ün yazdırılmadığını, yani bir toutine tamamlanmadan diğer routine’e geçildiğini de görebilirdik. Bunu da aşağıdaki örnekte gözlemleyelim.        
          
              
![concurrency vs paralellism](concurrency4.png)                     
          
Bu kez uygulamaya her bir metni yazdırdıktan sonra yarım saniye kadar beklemesini söyledik. Çıktıda da gördüğümüz üzere go scheduler go routine’leri öyle şekilde düzenledi ki; sıradan bağımsız ve zaman açısından oldukça efektif bir sonuç elde ettik.     

![concurrency vs paralellism](concurrency5.png)                 
            
       
  ### WAITGROUPS
İlk örnekte karşılaştığımız problemi main fonksiyonunun sonuna 1 saniyelik sleep metodu yazarak çözmüştük hatırlarsanız. Ancak bu çok da tercih edilesi bir yol değil. Waitgroup veya Channellar kullanarak goroutine’ler tamamlanmadan uygulamanının kapanmamasını sağlayabiliriz. Waitgrouplar adından da anlaşılabileceği gibi, bir grup işin tamamlanıncaya kadar belli bir noktada beklenmesi prensibiyle çalışır. Basit bir örnekle kullanımını inceleyelim:        
               
![concurrency vs paralellism](concurrency6.png)                        
          
Bu örnekte bir döngü oluşturularak döngünün her adımında wg waitgroup’una yeni bir yükleme yapılıyor. Bu durumun sonucu olarak uygulama, toplam 5 adet Done() metodunu alana kadar Wait() metodunun bulunduğu satırdan sonrasını çalıştırmıyor. Bu sayede uygulamanın devam etmesi için main fonksiyonuna sleep metodu yazmak zorunda kalmadık. Waitgrouplar belli goroutine’lerin tamamlanmadan uygulamanın devam etmesini istemediğimiz durumlarda kullanılabilir.            
             
             
  ### CHANNELS (KANALLAR)
Gorotineler birbirleriyle veri iletişimi kurabilir ve bunu channellar aracılığıyla yaparlar. Yine basit bir örnek ile başlayalım.         

![concurrency vs paralellism](concurrency7.png)                       
            
            
Channel’lar da tıpkı diğer değişkenler gibi int, string, slice gibi tiplerde olabilirler. Channelların tanımlanması make() metodu aracılığıyla yapılır. Yukarıdaki örnekte int tipinde bir channel tanımladık. Ve channel’ımızı printSomeText goroutine’ine gönderdik. Channel’lardan değer okumak veya channel’lara değer atamak için <- işaretinden faydalanır. Burada okun başlangıcı nereyi gösteriyorsa oradan, okun ucuna doğru veri aktarımı yapılır. Örneğimizde printSomeText fonksiyonu içerisinde de channel <- 5 kodu ile channel’a 5 değeri atanır. Bu kısım önemli, çünkü bu adım atlanırsa yani channel’a herhangi değer ataması gerçekleşmezse channel’ın gönderildiği rutine tekrar dönüp channel’ı okumak istediğimizde ‘deadlock’ hatası ile karşılaşırız. 
 ‘ received1 := <- chan1 ‘ satırının bir başka önemli noktası da channel’dan received1 değişkenine aktarım yapılana kadar beklenecek oluşudur. Aşağıdaki sonuçta görüldüğü üzere text 1 metnini yazdıktan sonra chan1 channel’ı henüz geri değer döndürmediği için 1 saniye bekledi. Channel kullanmasaydık bu şekilde bir durumla karşılaşmazdık. Ve yine aynı mantıkla main fonksiyonuna sleep metodunu eklemememize rağmen channel kullandığımız için ilk örneğin aksine ekranda hiçbir şey yazdırılmaması problemi ile karşılaşmadık.           
         
![concurrency vs paralellism](concurrency8.png)                 
        
