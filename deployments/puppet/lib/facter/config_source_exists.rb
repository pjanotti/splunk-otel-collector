# This custom fact checks for the existence of a collector configuration source file.
#
Facter.add(:collector_config_source_exists) do
    setcode do
        File.exist?($splunk_otel_collector::params::collector_config_source) ? true : false
    end
end
