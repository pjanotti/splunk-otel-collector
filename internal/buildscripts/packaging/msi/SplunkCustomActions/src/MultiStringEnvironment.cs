// Copyright  Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

using Microsoft.Win32;

public class MultiStringEnvironment : IDisposable
{
    private readonly string _valueName;
    private readonly RegistryKey _key;
    private readonly RegistryKey _subkey;

    public MultiStringEnvironment() :
        this(RegistryHive.LocalMachine, @"SYSTEM\CurrentControlSet\Services\splunk-otel-collector", "Environment")
    {
    }

    public MultiStringEnvironment(RegistryHive registryHive, string subKey, string valueName)
    {
        _valueName = valueName;
        _key = RegistryKey.OpenBaseKey(registryHive, RegistryView.Registry64);
        _subkey = _key.OpenSubKey(subKey, writable: true);
    }

    public string[] GetEnvironmentValue()
    {
        return (string[])(_subkey.GetValue(_valueName) ?? Array.Empty<string>());
    }

    public void AddEnvironmentVariables(Dictionary<string, string> environmentVariables)
    {
        string[] existingEnvironmentVariables = GetEnvironmentValue();
        string[] newEnvironment = 
            [.. existingEnvironmentVariables, .. environmentVariables.Select(kvp => $"{kvp.Key}={kvp.Value}")];

        // Sort the environment variables to ensure that the order is consistent
        Array.Sort(newEnvironment, StringComparer.OrdinalIgnoreCase);

        _subkey.SetValue(_valueName, newEnvironment, RegistryValueKind.MultiString);
    }

    public void Dispose()
    {
        _subkey.Close();
        _key.Close();
    }
}
