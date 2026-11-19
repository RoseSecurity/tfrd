class Tfrd < Formula
  desc "Find unused attributes and parameters in Terraform resources"
  homepage "https://github.com/RoseSecurity/tfrd"
  url "https://github.com/RoseSecurity/tfrd/archive/refs/tags/v0.3.1.tar.gz"
  sha256 "ad0077b1c858b9c3f60b5580f7e2e857da920db1685a558b5113f34fe52200f7"
  license "Apache-2.0"

  depends_on "go" => :build

  def install
    ldflags = "-s -w -X github.com/RoseSecurity/tfrd/cmd.Version=#{version}"
    system "go", "build", "-ldflags", ldflags, "-o", bin/"tfrd", "."
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/tfrd --version")
  end
end
