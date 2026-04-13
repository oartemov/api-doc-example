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

### Карточка

Замечание ИИ | Тип           | «Светофор»   | Приоритет    | Что делаем
------------- | -------------| -------------| -------------| -------------
Командно-инъекционная уязвимость, риск полной компроментации сервера|корректное| Accept| blocker| Исключить вызов shell, заменить на sybprocess.run(). Изучить возможность использования tarfile
Валидируйте `folder`: разрешите только безопасные символы (`[a-zA-Z0-9_./-]`) и проверяйте, что путь существует| Accept|major | Добавить проверку параметра folder на безопасные символы


## Пример 2 (Java): риск NullPointerException и неявная логика ошибок

### Промпт 

Ты - senior java developer.
Проведи ревью кода на Java:

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

Сервис отправляет письмо пользователю. 

Проанализируй код на предмет потенциальных багов:
- Обращение к свойствам без проверки на null/undefined
- Выход за границы массива
- Необработанные пустые строки

Формат ответа:
нумерованный список, для каждого бага (уязвимости) укажи: 
- blocker / major / minor
- что не так
- риск
- как исправить

### Ответ ИИ

>Вот ревью кода с нумерованным списком найденных проблем:
>
>---
>
>### 1. Отсутствие проверки `user` на `null`
>
>- **Уровень:** blocker  
>- **Что не так:** `userRepository.findById(userId)` может вернуть `null`, если пользователь с указанным `userId` не найден в базе. В следующей строке происходит вызов `user.getEmail()`, что гарантированно бросит `NullPointerException`.  
>- **Риск:** Падение приложения (необработанное исключение) при попытке отправить письмо несуществующему пользователю. Это приводит к ошибке в рантайме и прерыванию выполнения потока.  
>- **Как исправить:**  
>  ```java
>  User user = userRepository.findById(userId);
>  if (user == null) {
>      throw new UserNotFoundException("User not found: " + userId);
>  }
>  ```
>
>---
>
>### 2. Отсутствие проверки `userId` на `null`
>
>- **Уровень:** major  
>- **Что не так:** Метод принимает `Long userId`, но не проверяет, что аргумент не `null`. Если вызвать `sendWelcomeEmail(null)`, то `findById(null)` может либо выбросить NPE (в зависимости от реализации репозитория), либо вернуть `null`, что приведёт к проблеме из пункта 1.  
>- **Риск:** Неявное падение с `NullPointerException` без понятного сообщения об ошибке, затрудняющее отладку. Нарушение принципа явной проверки входных данных (fail-fast).  
>- **Как исправить:**  
>  ```java
>  if (userId == null) {
>      throw new IllegalArgumentException("userId must not be null");
>  }
>  ```
>
>---
>
>### 3. Отсутствие проверки `email` на `null`
>
>- **Уровень:** blocker  
>- **Что не так:** Даже если пользователь существует, его `email` может быть `null` (например, поле в БД допускает NULL). Вызов `mailClient.send(null, "Welcome!")` приведёт к `NullPointerException` внутри почтового клиента.  
>- **Риск:** Аналогично пункту 1 — неконтролируемое падение приложения. Дополнительно: если `mailClient` обработает `null` как строку "null", это приведёт к отправке письма на невалидный адрес (потенциальная проблема с доставляемостью).  
>- **Как исправить:**  
>  ```java
>  String email = user.getEmail();
>  if (email == null) {
>      throw new IllegalStateException("User " + userId + " has no email address");
>  }
>  mailClient.send(email, "Welcome!");
>  ```
>
>---
>
>### 4. Отсутствие проверки `email` на пустую строку
>
>- **Уровень:** major  
>- **Что не так:** Поле `email` может быть не `null`, но пустой строкой `""` или состоять только из пробелов `"   "`. Это не вызовет NPE, но приведёт к отправке письма на некорректный адрес (почтовый сервер, скорее всего, вернёт ошибку, либо письмо уйдёт в /dev/null).  
>- **Риск:** Тихая потеря функциональности — пользователь не получит письмо, но система не сообщит об ошибке. Это может остаться незамеченным, особенно при отсутствии мониторинга исходящих писем.  
>- **Как исправить:**  
>  ```java
>  if (email == null || email.trim().isEmpty()) {
>      throw new IllegalStateException("User " + userId + " has invalid (empty) email");
>  }
>  ```
>
>---
>
>### 5. Отсутствие обработки ошибок `mailClient.send()`
>
>- **Уровень:** minor (погранично с major)  
>- **Что не так:** Метод `send()` может выбросить исключение (сеть недоступна, SMTP-сервер не отвечает, авторизация не пройдена и т.д.), но текущий код никак его не обрабатывает. Исключение уйдёт наверх, и вызывающий код не поймёт, что письмо не отправлено, а пользователь может остаться без welcome-письма.  
>- **Риск:** Отсутствие возможности переотправки или логирования ошибок. Пользовательский опыт ухудшается без явного уведомления администратора. В зависимости от контекста, ошибка отправки письма может быть критичной (например, в сервисе подтверждения email).  
>- **Как исправить:**  
>  ```java
>  try {
>      mailClient.send(email, "Welcome!");
>  } catch (MailException e) {
>      log.error("Failed to send welcome email to user {}", userId, e);
>      throw new EmailSendingException("Cannot send welcome email", e);
>  }
>  ```
>
>---
>
>### Итоговая рекомендация по исправлению (полный безопасный код)
>
>```java
>import org.slf4j.Logger;
>import org.slf4j.LoggerFactory;
>
>public class MailService {
>    private static final Logger log = LoggerFactory.getLogger(MailService.class);
>    private final UserRepository userRepository;
>    private final MailClient mailClient;
>
>    public MailService(UserRepository userRepository, MailClient mailClient) {
>        this.userRepository = userRepository;
>        this.mailClient = mailClient;
>    }
>
>    public void sendWelcomeEmail(Long userId) {
>        if (userId == null) {
>            throw new IllegalArgumentException("userId must not be null");
>        }
>
>        User user = userRepository.findById(userId);
>        if (user == null) {
>            throw new UserNotFoundException("User not found: " + userId);
>        }
>
>        String email = user.getEmail();
>        if (email == null || email.trim().isEmpty()) {
>            throw new IllegalStateException("User " + userId + " has no valid email");
>        }
>
>        try {
>            mailClient.send(email, "Welcome!");
>            log.info("Welcome email sent to user {}", userId);
>        } catch (MailException e) {
>            log.error("Failed to send welcome email to user {}", userId, e);
>            throw new EmailSendingException("Cannot send welcome email to user " + userId, e);
>        }
>    }
>}
>```

### Карточка

Замечание ИИ | Тип           | «Светофор»   | Приоритет    | Что делаем
------------- | -------------| -------------| -------------| -------------
Отсутствие проверки `user` на `null`|корректное| Accept| blocker| Добавить проверку `user` на `null`
Отсутствие проверки `userId` на `null`|корректное| Accept| major| Добавить проверку `userId` на `null`
Отсутствие проверки `email` на `null`|корректное| Accept| blocker| Добавить проверку `email` на `null`
Отсутствие проверки `email` на пустую строку|корректное| Accept| major| Добавить проверку `email` на пустую строку

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
