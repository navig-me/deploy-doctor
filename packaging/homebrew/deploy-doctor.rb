class DeployDoctor < Formula
  desc "Diagnose container deploy readiness"
  homepage "https://github.com/navig-me/docker-doctor"
  version "0.1.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/navig-me/docker-doctor/releases/download/v0.1.0/deploy-doctor_0.1.0_darwin_arm64.tar.gz"
      sha256 "REPLACE_ME"
    else
      url "https://github.com/navig-me/docker-doctor/releases/download/v0.1.0/deploy-doctor_0.1.0_darwin_amd64.tar.gz"
      sha256 "REPLACE_ME"
    end
  end

  def install
    bin.install "deploy-doctor"
  end
end
