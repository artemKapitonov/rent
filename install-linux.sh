# Проверяем, установлен ли Go
if ! command -v go &> /dev/null; then
    echo "Go не установлен. Устанавливаем Go с помощью"
    curl -O https://dl.google.com/go/go1.22.1.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.22.1.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    export GOPATH=$HOME/go
    echo "export PATH=$PATH:$GOPATH/bin" >> ~/.bashrc
    echo "export PATH=$PATH:$GOPATH/bin" >> ~/.zshrc
fi

# Установка пакета
echo "Устанавливаем пакет rent..."

git clone https://github.com/artemKapitonov/rent
cd rent
go build
sudo mv rent /usr/local/bin

# Создаем директорию .config/rent
mkdir -p "$HOME/.config/rent"

# Создаем файл setting.yaml
touch "$HOME/.config/rent/setting.yaml"

# Проверяем, какой shell используется и добавляем команду автодополнения
if [[ $SHELL == *"zsh"* ]]; then
    echo "Добавляем автодополнение в .zshrc..."
    echo "rent completion zsh > /tmp/completion" >> ~/.zshrc
    echo "source /tmp/completion" >> ~/.zshrc
elif [[ $SHELL == *"bash"* ]]; then
    echo "Добавляем автодополнение в .bashrc..."
    echo "rent completion zsh > /tmp/completion" >> ~/.bashrc
    echo "source /tmp/completion" >> ~/.bashrc
else
    echo "Неизвестный shell: $SHELL. Автодополнение не будет добавлено."
fi

echo "Скрипт завершен. Пожалуйста, перезапустите терминал или выполните 'source ~/.bashrc' или 'source ~/.zshrc' для применения изменений."
