// swift-tools-version: 6.3
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
    name: "BrewBotPackages",
    platforms: [
            .iOS(.v17)
    ],
    products: [
        .library(name: "Domain", targets: ["Domain"]),
        .library(name: "Core", targets: ["Core"]),
        .library(name: "Data", targets: ["Data"]),
        .library(name: "Presentation", targets: ["Presentation"]),
    ],
    dependencies: [],
    targets: [
        .target(name: "Domain", dependencies: []),
        .target(name: "Core", dependencies: []),
        .target(name: "Data", dependencies: ["Domain", "Core"]),
        .target(name: "Presentation", dependencies: ["Domain", "Data", "Core"])
    ],
    swiftLanguageModes: [.v6]
)
