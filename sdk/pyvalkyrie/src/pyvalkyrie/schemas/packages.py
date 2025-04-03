from pydantic import BaseModel, Field

from .base import ResponseBase


class InstallNixPackage(BaseModel):
    msgType: str = Field(default="InstallNixPackage", frozen=True)
    pkgName: str = Field(description="Nix package to install")
    channel: str = Field(default="nixpkgs", description="Name of the channel")


class InstallNixPackageResponse(ResponseBase):
    msgType: str = Field(default="InstallNixPackageResponse", frozen=True)


class UninstallNixPackage(BaseModel):
    msgType: str = Field(default="UninstallNixPackage", frozen=True)
    pkgName: str = Field(description="Nix package to uninstall")


class UninstallNixPackageResponse(ResponseBase):
    msgType: str = Field(default="UninstallNixPackageResponse", frozen=True)
