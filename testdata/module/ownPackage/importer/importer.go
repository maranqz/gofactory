package importer

import "factory/ownPackage"

// ProducerShapedFuncIsReported: -ownPackage only restricts the owner
// package, so a function outside it that returns ownPackage.Loan — a
// producer shape inside the owner package — still must use the factory.
func ProducerShapedFuncIsReported() ownPackage.Loan {
	return ownPackage.Loan{} // want `Use factory for ownPackage.Loan`
}

// ConstConversionIsReported: the owner package's const exemption does not
// reach an importer's own const declaration.
const ConstConversionIsReported = ownPackage.Status(2) // want `Use factory for ownPackage.Status`
