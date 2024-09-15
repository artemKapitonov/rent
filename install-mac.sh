# Проверяем, установлен ли Go
if ! command -v go &> /dev/null; then
    echo "Go не установлен. Устанавливаем Go с помощью Homebrew..."
    # Установка Go через Homebrew
    brew install go
fi

# Установка пакета
echo "Устанавливаем пакет rent..."
go install github.com/artemKapitonov/rent

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
