metadata    :name        => "dht220",
            :description => "Interact with a DHT220 via ChoriaPI",
            :author      => "R.I.Pienaar <rip@devco.net>",
            :license     => "Apache License, Version 2.0",
            :version     => "0.0.1",
            :url         => "https://www.devco.net/",
            :timeout     => 10

action "reading", :description => "Gather temp and humidity" do
    display :always

    output "temperature",
        :description => "The current temperature",
        :display_as => "Temperature"

    output "humidity",
        :description => "The current humidity",
        :display_as => "Humidity"    

    output "time",
        :description => "Measurement time",
        :display_as => "Time"
end