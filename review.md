# Ревью кода + классификация ошибок

## Пример 1 (Python): команда из пользовательского ввода (риск безопасности)

Ты - инженер по безопасности.
Проведи ревью кода на Python:

```python
import os
def run_backup(folder: str) -> int:    # folder приходит извне (например, из формы)     
cmd = f"tar -czf backup.tgz {folder}"    
    return os.system(cmd)
```
	
Утилита запускает системную команду. Аргумент приходит от пользователя. 
Проверь его по списку OWASP Top 10.

Формат ответа:
- blocker / major / minor
- что не так
- риск
- как исправить

### Ответ ИИ

>**Ответ: blocker**

>**Что не так:**  
>Команда формируется через f-строку напрямую с пользовательским вводом (`folder`). Это приводит к **командно-инъекционной уязвимости (Command Injection)**. Злоумышленник может передать значение вроде `"docs; rm -rf /"` или `"$(malicious_command)"`, и ОС выполнит произвольные команды.
>
>**Риск:**  
>- Полная компрометация сервера (удаление/модификация файлов, кража данных, установка бэкдоров).  
>- Соответствует **A03:2021 – Injection** (и классической **A1:2017 – Injection**) в OWASP Top 10.  
>- В контексте Python — критический уровень опасности.
>
>**Как исправить:**  
>
>1. **Не используйте `os.system()`** — она вызывает shell и уязвима по умолчанию.  
>2. **Используйте `subprocess.run()` с массивом аргументов** и без `shell=True`:  
>   ```python
>   import subprocess
>
>   def run_backup(folder: str) -> int:
>       # Проверка, что folder не содержит опасных символов (дополнительно)
>       if not folder.isalnum() and not all(c.isalnum() or c in '/._-' for c in folder):
>           raise ValueError("Недопустимые символы в имени папки")
>       
>       # Вызов без shell
>       result = subprocess.run(['tar', '-czf', 'backup.tgz', folder])
>       return result.returncode
>   ```
>3. **Ещё лучше — использовать `shlex.quote()`** (но только если shell неизбежен, что не рекомендуется):  
>   ```python
>   import shlex
>   cmd = f"tar -czf backup.tgz {shlex.quote(folder)}"
>   ```
>4. **Идеально — заменить tar на Python-библиотеку** (например, `tarfile`), полностью исключив вызов системной команды:  
>   ```python
>   import tarfile
>
>   def run_backup(folder: str) -> int:
>       with tarfile.open('backup.tgz', 'w:gz') as tar:
>           tar.add(folder, arcname=os.path.basename(folder))
>       return 0  # или реальный код ошибки
>   ```
>
>**Дополнительно:**  
>- Валидируйте `folder`: разрешите только безопасные символы (`[a-zA-Z0-9_./-]`) и проверяйте, что путь существует.  
>- Не полагайтесь только на экранирование — лучшая защита **не вызывать shell** для пользовательских данных.

# Карточка

Замечание ИИ | Тип           | «Светофор»   | Приоритет    | Что делаем
------------- | -------------| -------------| -------------| -------------
Командно-инъекционная уязвимость, риск полной компроментации сервера|корректное| Accept| blocker| Исключить вызов shell, заменить на sybprocess.run(). Изучить возможность использования tarfile

· : …
· Тип: 
· : Accept / Verify / Reject
· Приоритет: blocker / major / minor
· Что делаем дальше: …

## Пример 2 (Java): риск NullPointerException и неявная логика ошибок

```public class MailService {  
private final UserRepository userRepository;  
private final MailClient mailClient;  
public MailService(UserRepository userRepository, MailClient mailClient) {    
    this.userRepository = userRepository;    
    this.mailClient = mailClient;  
}  
    public void sendWelcomeEmail(Long userId) {    
        User user = userRepository.findById(userId);     
				mailClient.send(user.getEmail(), "Welcome!");  
    }
}```

Сервис отправляет письмо пользователю. Репозиторий может вернуть null.

## Пример 3 (Go): SQL-инъекция через конкатенацию строки

Проверь код: 

```package handlers
import (
“database/sql”
“fmt”
“net/http”
)

func Orders(db *sql.DB, w http.ResponseWriter, r *http.Request) error {
status := r.URL.Query().Get(“status”)
query := fmt.Sprintf(“SELECT id, total FROM orders WHERE status = ‘%s’”, status)
_, err := db.Query(query)
return err
}```

Обработчик принимает параметр status и собирает SQL строкой.
Важно:
• найди уязвимости (SQL injection);
• оцени приоритет (blocker/major/minor);
• предложи минимально безопасное исправление
Формат ответа:
- blocker / major / minor
- что не так
- риск
- как исправить
