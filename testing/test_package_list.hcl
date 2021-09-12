package "test" {
    base_dir="/base"

    version "latest" {
        image="packageless/testing"
    
        volume {
            path="/a/path"
            mount="/mount/path"
        }

        volume {
            mount="/run/"
        }

        copy {
            source="/a/source"
            dest="/a/dest"
        }

        port="3000"
    }
    
}