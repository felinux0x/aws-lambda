function Run-Build {
    Write-Host "🚀 Building for AWS Lambda (Linux/ARM64)..." -ForegroundColor Cyan
    # Salvando estado anterior para não quebrar o terminal do usuário
    $oldGoos = $env:GOOS
    $oldGoarch = $env:GOARCH
    $oldCgo = $env:CGO_ENABLED

    $env:GOOS="linux"
    $env:GOARCH="arm64"
    $env:CGO_ENABLED="0"
    
    go build -tags lambda.norpc -o bootstrap cmd/api/main.go
    
    # Restaurando estado original
    $env:GOOS = $oldGoos
    $env:GOARCH = $oldGoarch
    $env:CGO_ENABLED = $oldCgo
    
    Write-Host "✅ Build complete: ./bootstrap" -ForegroundColor Green
}

function Run-Test {
    Write-Host "🧪 Running Tests..." -ForegroundColor Cyan
    # Resetando GOOS para garantir que o teste rode no Windows
    $env:GOOS = $null
    $env:GOARCH = $null
    go test -v ./internal/service/... ./internal/handler/...
}

function Run-Bench {
    Write-Host "📊 Running Benchmarks..." -ForegroundColor Cyan
    $env:GOOS = $null
    $env:GOARCH = $null
    # Comando corrigido para focar no pacote de service
    go test -v -bench=. -benchmem github.com/feliperosa/aws-lambda-go/internal/service
}

function Run-Local {
    if (Get-Command sam -ErrorAction SilentlyContinue) {
        Run-Build
        Write-Host "🏠 Invoking SAM Local..." -ForegroundColor Cyan
        sam local invoke ApiFunction --event events/api_event.json --env-vars env.json
    } else {
        Write-Host "❌ Error: AWS SAM CLI not found." -ForegroundColor Red
        Write-Host "Please install it from: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html" -ForegroundColor Yellow
    }
}

function Run-Debug {
    if (Get-Command sam -ErrorAction SilentlyContinue) {
        Write-Host "🚧 Building for Debugging..." -ForegroundColor Cyan
        $oldGoos = $env:GOOS
        $env:GOOS="linux"
        $env:GOARCH="arm64"
        $env:CGO_ENABLED="0"
        go build -gcflags="all=-N -l" -o bootstrap cmd/api/main.go
        $env:GOOS = $oldGoos
        
        Write-Host "🐞 Starting SAM in Debug Mode (Port 5984)..." -ForegroundColor Magenta
        sam local invoke ApiFunction --event events/api_event.json --env-vars env.json -d 5984 --debugger-path (Get-Command dlv).Source | Select-Object -ExpandProperty Parent
    } else {
        Write-Host "❌ Error: AWS SAM CLI not found." -ForegroundColor Red
    }
}

switch ($args[0]) {
    "build" { Run-Build }
    "test"  { Run-Test }
    "bench" { Run-Bench }
    "local" { Run-Local }
    "debug" { Run-Debug }
    default {
        Write-Host "Usage: ./run.ps1 [build|test|bench|local|debug]" -ForegroundColor Yellow
    }
}
