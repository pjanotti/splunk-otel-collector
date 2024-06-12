# Class for collecting the configuration options for the splunk-otel-collector service
class splunk_otel_collector::collector_win_config_options {
    $base_env_vars = {
        'SPLUNK_ACCESS_TOKEN' => $splunk_otel_collector::splunk_access_token,
        'SPLUNK_API_URL' => $splunk_otel_collector::splunk_api_url,
        'SPLUNK_BUNDLE_DIR' => $splunk_otel_collector::splunk_bundle_dir,
        'SPLUNK_COLLECTD_DIR' => $splunk_otel_collector::splunk_collectd_dir,
        'SPLUNK_CONFIG' => $splunk_otel_collector::collector_config_dest,
        'SPLUNK_HEC_TOKEN' => $splunk_otel_collector::splunk_hec_token,
        'SPLUNK_HEC_URL' => $splunk_otel_collector::splunk_hec_url,
        'SPLUNK_INGEST_URL' => $splunk_otel_collector::splunk_ingest_url,
        'SPLUNK_MEMORY_TOTAL_MIB' => $splunk_otel_collector::splunk_memory_total_mib,
        'SPLUNK_REALM' => $splunk_otel_collector::splunk_realm,
        'SPLUNK_TRACE_URL' => $splunk_otel_collector::splunk_trace_url,
    }

    $ballast_size_mib = $splunk_otel_collector::collector_version != 'latest' and
        versioncmp($splunk_otel_collector::collector_version, '0.97.0') < 0 and
        !$splunk_otel_collector::splunk_ballast_size_mib.strip().empty() ? {
            true    => { 'SPLUNK_BALLAST_SIZE_MIB' => $splunk_otel_collector::splunk_ballast_size_mib },
            default => {},
        }

    $gomemlimit = ($splunk_otel_collector::collector_version == 'latest' or
        versioncmp($splunk_otel_collector::collector_version, '0.97.0') >= 0) and
        !$splunk_otel_collector::gomemlimit.strip().empty() ? {
            true    => { 'GOMEMLIMIT' => $splunk_otel_collector::gomemlimit },
            default => {},
        }

    $listen_interface = !$splunk_otel_collector::splunk_listen_interface.strip().empty() ? {
        true    => { 'SPLUNK_LISTEN_INTERFACE' => $splunk_otel_collector::splunk_listen_interface },
        default => {},
    }

    $collector_env_vars = merge($base_env_vars, $ballast_size_mib, $gomemlimit, $listen_interface)
}