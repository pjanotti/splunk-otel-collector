# escape=`

ARG windowscontainer_version=ltsc2019
FROM mcr.microsoft.com/dotnet/framework/aspnet:4.8-windowsservercore-${windowscontainer_version}

COPY apps/bin/_PublishedWebsites/AspNet.WebApi.NetFramework/ /apps/aspnetfx

# Add test application to the Default Web Site
# Escape double-quotes since the command is passed as a parameter to PowerShell.exe
RUN /Windows/System32/inetsrv/appcmd add app /site.name:\"Default Web Site\" /path:/aspnetfxapp /physicalPath:\"c:\apps\aspnetfx\"
