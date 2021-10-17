pim "python" {
    base_dir="/python/"
   
    version "latest" {
        image="packageless/python"
        
        volume {
            path="/python/packages/"
            mount="/usr/local/lib/python3.9/site-packages/"
        }

        volume {
            mount="/run/"
        }

        copy {
            source="/usr/local/lib/python3.9/site-packages/"
            dest="/python/packages/"
        }

        port="3000"
    }
}