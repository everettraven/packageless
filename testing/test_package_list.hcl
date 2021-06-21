package "test" {
    image="packageless/testing"
    base_dir="/base"
    
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