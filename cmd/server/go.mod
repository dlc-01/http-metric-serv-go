module server

go 1.20

require storage v1.0.0
replace storage v1.0.0 => ../../internal/storage
require handlers v1.0.0
replace handlers v1.0.0 => ../../internal/handlers