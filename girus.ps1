function Invoke-Girus {
    param(
        [Parameter(Position=0, ValueFromRemainingArguments=$true)]
        [string[]]$Arguments
    )
    
    $cmd = "docker run --rm -v /var/run/docker.sock:/var/run/docker.sock"
    if ($env:USERPROFILE) {
        $cmd += " -v `"$env:USERPROFILE/.girus:/root/.girus`""
    }
    $cmd += " girus $($Arguments -join ' ')"
    
    Invoke-Expression $cmd
}

# Criar alias para facilitar o uso
Set-Alias -Name girus -Value Invoke-Girus

# Exportar a função para que esteja disponível em outras sessões
Export-ModuleMember -Function Invoke-Girus -Alias girus