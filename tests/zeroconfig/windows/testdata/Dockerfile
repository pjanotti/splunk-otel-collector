# escape=`

ARG windowscontainer_version=ltsc2022
FROM mcr.microsoft.com/dotnet/framework/aspnet:4.8-windowsservercore-${windowscontainer_version}

# On ltsc2019 RUN defaults to "powershell -Command {0}", while ltsc2022 RUN defaults to "cmd /S /C"
# set the SHELL so it is consistent between both versions.
SHELL [ "powershell", "-Command" ]

COPY apps/bin/_PublishedWebsites/AspNet.WebApi.NetFramework/ /apps/aspnetfx

# Add test application to the Default Web Site
RUN \Windows\System32\inetsrv\appcmd add app /site.name:\"Default Web Site\" /path:/aspnetfxapp /physicalPath:\"c:\apps\aspnetfx\"
