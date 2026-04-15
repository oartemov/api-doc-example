# Ревью кода + классификация ошибок

## Пример 1 (Python): команда из пользовательского ввода (риск безопасности)

### Промпт 

Ты - инженер по безопасности.
Проведи ревью кода на Python на наличие уязвимостей:

```python
import os
def run_backup(folder: str) -> int:    # folder приходит извне (например, из формы)     
cmd = f"tar -czf backup.tgz {folder}"    
    return os.system(cmd)
```
	
Утилита запускает системную команду. Аргумент приходит от пользователя. 
Проверь его по списку OWASP Top 10.

Формат ответа:
нумерованный список, для каждой уязвимости укажи: 
- blocker / major / minor
- что не так
- риск
- как исправить

Сообщай только о реально существующих проблемах, не нужно описывать теоретические риски.


### Ответ ИИ

>Вот ревью кода на основе **OWASP Top 10** (актуальная версия).
>
>---
>
>### 1. A03:2021 – Injection (Command Injection)
>- **Severity**: **blocker**
>- **Что не так**: Пользовательский ввод (`folder`) напрямую подставляется в команду оболочки через f-строку. Злоумышленник может передать значение типа `"docs; rm -rf /"`, что выполнит произвольные системные команды.
>- **Риск**: Полная компрометация сервера, удаление данных, установка бэкдоров, кража конфиденциальной информации.
>- **Как исправить**:
>  - Использовать `subprocess.run()` с массивом аргументов (без `shell=True`).
>  - Экранировать путь с помощью `shlex.quote()` (но лучше — массив).
>  - Пример:
>    ```python
>    import subprocess
>    def run_backup(folder: str) -> int:
>        result = subprocess.run(["tar", "-czf", "backup.tgz", folder])
>        return result.returncode
>    ```
>
>---
>
>### 2. A04:2021 – Insecure Design
>- **Severity**: **major**
>- **Что не так**: Функция не проверяет, что `folder` является допустимым путём (например, не содержит `..` или ссылок на системные директории). Также нет ограничений на размер архива или время выполнения.
>- **Риск**: Злоумышленник может архивировать чувствительные системные файлы (`/etc/passwd`, `/root/.ssh`), выйти за пределы разрешённой директории, вызвать отказ в обслуживании (гигантский архив).
>- **Как исправить**:
>  - Белый список разрешённых базовых директорий.
>  - Нормализовать путь с `os.path.abspath()` и проверить, что он начинается с разрешённого префикса.
>  - Установить ограничение на размер архива и таймаут выполнения.
>
>---
>
>### 3. A08:2021 – Software and Data Integrity Failures
>- **Severity**: **minor**
>- **Что не так**: Файл `backup.tgz` создаётся в текущей рабочей директории, которая может быть предсказуема или перезаписываема другим пользователем. Нет проверки, что файл не был подменён в процессе создания.
>- **Риск**: Подмена резервной копии, подсовывание вредоносного архива, потеря целостности данных.
>- **Как исправить**:
>  - Создавать архив во временной защищённой директории (например, `tempfile.mkdtemp()`).
>  - Вычислять и проверять хеш (SHA-256) после создания архива.
>  - Использовать атомарное перемещение файла после успешного создания.
>
>---
>
>### 4. A09:2021 – Security Logging and Monitoring Failures
>- **Severity**: **minor**
>- **Что не так**: Нет логирования того, кто, откуда и с каким параметром запустил бэкап. При инциденте невозможно восстановить последовательность действий.
>- **Риск**: Атака может остаться незамеченной, затруднён forensic-анализ.
>- **Как исправить**:
>  - Добавить логирование с уровнем INFO: пользователь, IP (если есть), `folder`, результат выполнения.
>  - Фиксировать неудачные попытки (невалидный путь, ошибки tar) как WARNING/ERROR.
>  - Убедиться, что логи отправляются в централизованную систему.
>
>---
>
>### Дополнительно (вне OWASP Top 10, но важно):
>- **Ошибка**: в коде отсутствует `import os` в приведённом фрагменте (подразумевается, но не указан). В реальном ревью это `minor` стилистическое замечание.
>- **Риск**: Функция возвращает `int` от `os.system()`, но при ошибке `tar` код возврата может быть непредсказуемым — лучше использовать `subprocess` с проверкой `stderr`.

### Карточка

Замечание ИИ | Тип           | «Светофор»   | Приоритет    | Что делаем
------------- | -------------| -------------| -------------| -------------
Пользовательский ввод напрямую подставляется в команду оболочки через f-строку|корректное| accept| blocker| Использовать `subprocess.run()` с массивом аргументов (без `shell=True`).
Функция не проверяет, что `folder` является допустимым путём|корректное| accept| major| Проверка параметра `folder` на допустимость, ограничение размера архива и установка таймаута
Файл `backup.tgz` создаётся в текущей рабочей директории|корректное| accept| minor| Создавать архив во временной защищённой директории, проверять хэш
Нет логирования того, кто, откуда и с каким параметром запустил бэкап|корректное| accept| minor| Добавить логирование

### Комментарии

#### 1. Command Injection

- **Что не так**: Пользовательский ввод (`folder`) напрямую подставляется в команду оболочки через f-строку. Злоумышленник может передать значение типа `"docs; rm -rf /"`, что выполнит произвольные системные команды.
- **Риск**: Полная компрометация сервера, удаление данных, установка бэкдоров, кража конфиденциальной информации.
- **Как исправить**: Использовать `subprocess.run()` с массивом аргументов (без `shell=True`).
- **Приоритет**: blocker

#### 2. Insecure Design

- **Что не так**: Функция не проверяет, что `folder` является допустимым путём (например, не содержит `..` или ссылок на системные директории). Также нет ограничений на размер архива или время выполнения.
- **Риск**: Злоумышленник может архивировать чувствительные системные файлы (`/etc/passwd`, `/root/.ssh`), выйти за пределы разрешённой директории, вызвать отказ в обслуживании (гигантский архив).
- **Как исправить**:
  - Белый список разрешённых базовых директорий.
  - Нормализовать путь с `os.path.abspath()` и проверить, что он начинается с разрешённого префикса.
  - Установить ограничение на размер архива и таймаут выполнения.
- **Приоритет**: major

#### 3. Software and Data Integrity Failures

- **Что не так**: Файл `backup.tgz` создаётся в текущей рабочей директории, которая может быть предсказуема или перезаписываема другим пользователем. Нет проверки, что файл не был подменён в процессе создания.
- **Риск**: Подмена резервной копии, подсовывание вредоносного архива, потеря целостности данных.
- **Как исправить**:
  - Создавать архив во временной защищённой директории (например, `tempfile.mkdtemp()`).
  - Вычислять и проверять хеш (SHA-256) после создания архива.
  - Использовать атомарное перемещение файла после успешного создания.
- **Приоритет**: minor

#### 4. Security Logging and Monitoring Failures

- **Что не так**: Нет логирования того, кто, откуда и с каким параметром запустил бэкап. При инциденте невозможно восстановить последовательность действий.
- **Риск**: Атака может остаться незамеченной, затруднён forensic-анализ.
- **Как исправить**:
  - Добавить логирование с уровнем INFO: пользователь, IP (если есть), `folder`, результат выполнения.
  - Фиксировать неудачные попытки (невалидный путь, ошибки tar) как WARNING/ERROR.
  - Убедиться, что логи отправляются в централизованную систему.
- **Приоритет**: minor

### Next Steps

#### 1. Что исправляем ДО MERGE (обязательные изменения)

##### Блокер (Command Injection) — исправить немедленно
```python
# Вместо
subprocess.run(f"tar -czf backup.tgz {folder}", shell=True)

# Сделать
subprocess.run(["tar", "-czf", "backup.tgz", folder], check=False)
```
- Полностью убрать `shell=True`
- Передать аргументы массивом

##### major (Insecure Design) — минимальная защита
- Добавить **базовую валидацию пути**:
  ```python
  import os
  
  ALLOWED_BASE_DIRS = ["/var/backups/allowed", "/home/user/data"]
  
  real_path = os.path.realpath(folder)
  if not any(real_path.startswith(base) for base in ALLOWED_BASE_DIRS):
      raise ValueError(f"Path {folder} is not allowed")
  ```
- Добавить **ограничение по времени**:
  ```python
  subprocess.run([...], timeout=300)  # 5 минут
  ```

#### 2. Что выносим в ОТДЕЛЬНЫЕ ЗАДАЧИ (следующий спринт)

##### Задача 1: Полная защита от path traversal (улучшение Insecure Design)
- Реализовать проверку на symlink-атаки (`os.path.realpath` уже частично решает)
- Добавить проверку, что целевая директория действительно существует и доступна для чтения
- Реализовать ограничение на **размер результирующего архива** (проверять после создания, >1GB — reject)

##### Задача 2: Целостность данных (Software and Data Integrity)
- Создавать архив во временной директории через `tempfile.mkdtemp()`
- Вычислять SHA-256 хеш после создания
- Атомарно перемещать файл в финальное место с проверкой хеша

##### Задача 3: Безопасное логирование (Security Logging)
- Добавить структурированное логирование (JSON-формат)
- Фиксировать: `user_id`, `ip_address`, `folder`, `archive_size`, `duration_ms`, `status`
- Настроить отправку логов в централизованную систему (ELK / Sentry / SIEM)

#### 3. Какие тесты добавить

##### Критические (до merge)
```python
def test_command_injection_prevention():
    """Проверка, что shell-инъекции не выполняются"""
    malicious_input = "docs; rm -rf /"
    with pytest.raises(FileNotFoundError):  # или другой ожидаемый failure
        backup_function(malicious_input)
    # Убедиться, что файловая система не повреждена
    assert os.path.exists("/")  # базовая проверка

def test_path_traversal_blocked():
    """Проверка блокировки .. и абсолютных путей вне allowed"""
    with pytest.raises(ValueError):
        backup_function("/etc/passwd")
    with pytest.raises(ValueError):
        backup_function("../sensitive")

def test_timeout_enforced():
    """Таймаут при архивации большого объема"""
    with pytest.raises(subprocess.TimeoutExpired):
        backup_function("/very/large/directory", timeout=1)
```

##### В отдельной задаче (после merge)
```python
def test_atomic_archive_creation():
    """Архив не появляется частично записанным"""
    # Проверить, что во время создания нет промежуточного файла в финальной локации

def test_checksum_verification():
    """Хеш созданного архива соответствует ожидаемому"""
    backup_function("/test/data")
    assert verify_checksum(archive_path) is True

def test_logging_contains_audit_info():
    """Логи содержат необходимые поля для аудита"""
    with capture_logs() as logs:
        backup_function("/allowed/folder")
    assert any("user_id" in log and "ip_address" in log for log in logs)

def test_concurrent_backups():
    """Параллельные бэкапы не пересекаются во временных файлах"""
    # Запустить несколько concurrent backup и проверить отсутствие race conditions
```

##### Тесты производительности (отдельный бенчмарк)
```python
def test_large_directory_timeout():
    """Проверка, что timeout действительно прерывает выполнение"""
    # Создать директорию с 10,000 файлов
    # Установить timeout=5
    # Проверить, что операция прервана
```

---


## Пример 2 (Java): риск NullPointerException и неявная логика ошибок

### Промпт 

Ты - senior java developer.
Проведи ревью кода на Java:

```java
public class MailService {  
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
}
```

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

Не предлагай крупный рефакторинг без причины, не придумывай новые библиотеки.

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
Отсутствие проверки `user` на `null`|корректное| accept| blocker| Добавить проверку `user` на `null`
Отсутствие проверки `userId` на `null`|корректное| accept| major| Добавить проверку `userId` на `null`
Отсутствие проверки `email` на `null`|корректное| accept| blocker| Добавить проверку `email` на `null`
Отсутствие проверки `email` на пустую строку|корректное| accept| major| Добавить проверку `email` на пустую строку
Отсутствие обработки ошибок `mailClient.send()`|корректное| accept| minor| Добавить обработку ошибок

### Комментарии

#### 1. Отсутствие проверки `user` на `null`

- **Что не так:** `userRepository.findById(userId)` может вернуть `null`, если пользователь с указанным `userId` не найден в базе. В следующей строке происходит вызов `user.getEmail()`, что гарантированно бросит `NullPointerException`.  
- **Риск:** Падение приложения (необработанное исключение) при попытке отправить письмо несуществующему пользователю. Это приводит к ошибке в рантайме и прерыванию выполнения потока.  
- **Как исправить:** добавить проверку на `null`
  ```java
  User user = userRepository.findById(userId);
  if (user == null) {
      throw new UserNotFoundException("User not found: " + userId);
  }
  ```
- **Приоритет**: blocker 

#### 2. Отсутствие проверки `email` на `null`

- **Что не так:** Даже если пользователь существует, его `email` может быть `null` (например, поле в БД допускает NULL). Вызов `mailClient.send(null, "Welcome!")` приведёт к `NullPointerException` внутри почтового клиента.  
- **Риск:** Аналогично пункту 1 — неконтролируемое падение приложения. Дополнительно: если `mailClient` обработает `null` как строку "null", это приведёт к отправке письма на невалидный адрес (потенциальная проблема с доставляемостью).  
- **Как исправить:** добавить проверку на `null`
  ```java
  String email = user.getEmail();
  if (email == null) {
      throw new IllegalStateException("User " + userId + " has no email address");
  }
  mailClient.send(email, "Welcome!");
  ```
- **Приоритет**: blocker

#### 3. Отсутствие проверки `userId` на `null`

- **Что не так:** Метод принимает `Long userId`, но не проверяет, что аргумент не `null`. Если вызвать `sendWelcomeEmail(null)`, то `findById(null)` может либо выбросить NPE (в зависимости от реализации репозитория), либо вернуть `null`, что приведёт к проблеме из пункта 1.  
- **Риск:** Неявное падение с `NullPointerException` без понятного сообщения об ошибке, затрудняющее отладку. Нарушение принципа явной проверки входных данных (fail-fast).  
- **Как исправить:** добавить проверку на `null`
  ```java
  if (userId == null) {
      throw new IllegalArgumentException("userId must not be null");
  }
  ```
- **Приоритет**: major

#### 4. Отсутствие проверки `email` на пустую строку

- **Что не так:** Поле `email` может быть не `null`, но пустой строкой `""` или состоять только из пробелов `"   "`. Это не вызовет NPE, но приведёт к отправке письма на некорректный адрес (почтовый сервер, скорее всего, вернёт ошибку, либо письмо уйдёт в /dev/null).  
- **Риск:** Тихая потеря функциональности — пользователь не получит письмо, но система не сообщит об ошибке. Это может остаться незамеченным, особенно при отсутствии мониторинга исходящих писем.  
- **Как исправить:** добавить проверку на пустую строку
  ```java
  if (email == null || email.trim().isEmpty()) {
      throw new IllegalStateException("User " + userId + " has invalid (empty) email");
  }
  ```
- **Приоритет**: major

#### 5. Отсутствие обработки ошибок `mailClient.send()`

- **Уровень:** minor (погранично с major)  
- **Что не так:** Метод `send()` может выбросить исключение (сеть недоступна, SMTP-сервер не отвечает, авторизация не пройдена и т.д.), но текущий код никак его не обрабатывает. Исключение уйдёт наверх, и вызывающий код не поймёт, что письмо не отправлено, а пользователь может остаться без welcome-письма.  
- **Риск:** Отсутствие возможности переотправки или логирования ошибок. Пользовательский опыт ухудшается без явного уведомления администратора. В зависимости от контекста, ошибка отправки письма может быть критичной (например, в сервисе подтверждения email).  
- **Как исправить:**  добавить обработку ошибок
  ```java
  try {
      mailClient.send(email, "Welcome!");
  } catch (MailException e) {
      log.error("Failed to send welcome email to user {}", userId, e);
      throw new EmailSendingException("Cannot send welcome email", e);
  }
  ```
- **Приоритет**: minor  

### Next Steps

#### 1. Что исправляем ДО MERGE (обязательные изменения)

##### Блокеры (Blocker) — исправить немедленно

**Пункт 1: Проверка user на null**
```java
User user = userRepository.findById(userId);
if (user == null) {
    throw new UserNotFoundException("User not found: " + userId);
}
```

**Пункт 2: Проверка email на null**
```java
String email = user.getEmail();
if (email == null) {
    throw new IllegalStateException("User " + userId + " has no email address");
}
mailClient.send(email, "Welcome!");
```

##### Major (обязательно до merge)

**Пункт 3: Проверка userId на null**
```java
public void sendWelcomeEmail(Long userId) {
    if (userId == null) {
        throw new IllegalArgumentException("userId must not be null");
    }
    // остальной код
}
```

**Пункт 4: Проверка email на пустую строку**
```java
if (email == null || email.trim().isEmpty()) {
    throw new IllegalStateException("User " + userId + " has invalid (empty) email");
}
```

#### 2. Что выносим в ОТДЕЛЬНЫЕ ЗАДАЧИ (следующий спринт)

##### Задача 1: Обработка ошибок mailClient (Minor - пункт 5)

**Требования:**
- Добавить try-catch блок для `mailClient.send()`
- Логировать ошибку с уровнем ERROR
- Бросать кастомное исключение `EmailSendingException`

#### 3. Какие тесты добавить

##### Критические тесты (до merge)

```java
@Test
void shouldThrowExceptionWhenUserIdIsNull() {
    // Given
    Long nullUserId = null;
    
    // When & Then
    assertThatThrownBy(() -> welcomeEmailService.sendWelcomeEmail(nullUserId))
        .isInstanceOf(IllegalArgumentException.class)
        .hasMessage("userId must not be null");
}

@Test
void shouldThrowExceptionWhenUserNotFound() {
    // Given
    Long nonExistentUserId = 999L;
    when(userRepository.findById(nonExistentUserId)).thenReturn(null);
    
    // When & Then
    assertThatThrownBy(() -> welcomeEmailService.sendWelcomeEmail(nonExistentUserId))
        .isInstanceOf(UserNotFoundException.class)
        .hasMessage("User not found: " + nonExistentUserId);
}

@Test
void shouldThrowExceptionWhenUserEmailIsNull() {
    // Given
    Long userId = 1L;
    User user = new User();
    user.setEmail(null);
    when(userRepository.findById(userId)).thenReturn(user);
    
    // When & Then
    assertThatThrownBy(() -> welcomeEmailService.sendWelcomeEmail(userId))
        .isInstanceOf(IllegalStateException.class)
        .hasMessageContaining("no email address");
}

@Test
void shouldThrowExceptionWhenUserEmailIsEmpty() {
    // Given
    Long userId = 1L;
    User user = new User();
    user.setEmail("");
    when(userRepository.findById(userId)).thenReturn(user);
    
    // When & Then
    assertThatThrownBy(() -> welcomeEmailService.sendWelcomeEmail(userId))
        .isInstanceOf(IllegalStateException.class)
        .hasMessageContaining("invalid email");
}

@Test
void shouldThrowExceptionWhenUserEmailIsBlank() {
    // Given
    Long userId = 1L;
    User user = new User();
    user.setEmail("   ");
    when(userRepository.findById(userId)).thenReturn(user);
    
    // When & Then
    assertThatThrownBy(() -> welcomeEmailService.sendWelcomeEmail(userId))
        .isInstanceOf(IllegalStateException.class)
        .hasMessageContaining("invalid email");
}

@Test
void shouldSendEmailSuccessfullyWhenAllValid() {
    // Given
    Long userId = 1L;
    String validEmail = "user@example.com";
    User user = new User();
    user.setEmail(validEmail);
    when(userRepository.findById(userId)).thenReturn(user);
    
    // When
    welcomeEmailService.sendWelcomeEmail(userId);
    
    // Then
    verify(mailClient).send(validEmail, "Welcome!");
}
```

##### Тесты для отдельной задачи (после добавления error handling)

```java
@Test
void shouldHandleMailClientException() {
    // Given
    Long userId = 1L;
    User user = new User();
    user.setEmail("user@example.com");
    when(userRepository.findById(userId)).thenReturn(user);
    doThrow(new MailException("SMTP error")).when(mailClient).send(any(), any());
    
    // When & Then
    assertThatThrownBy(() -> welcomeEmailService.sendWelcomeEmail(userId))
        .isInstanceOf(EmailSendingException.class)
        .hasCauseInstanceOf(MailException.class);
    
    // Verify logging (если используете spy на логгер)
    verify(log).error(contains("Failed to send welcome email"), eq(userId), any());
}

@Test
void shouldLogSuccessWhenEmailSent() {
    // Given
    Long userId = 1L;
    User user = new User();
    user.setEmail("user@example.com");
    when(userRepository.findById(userId)).thenReturn(user);
    
    // When
    welcomeEmailService.sendWelcomeEmail(userId);
    
    // Then
    verify(log).info(contains("sent successfully"), eq(userId));
}
```

---

## Пример 3 (Go): SQL-инъекция через конкатенацию строки

### Промпт 

Ты - инженер по безопасности.

Проверь код на go на наличие уязвимостей: 

```go
package handlers
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
}
```

Обработчик принимает параметр status и собирает SQL строкой.
Важно:
- найди уязвимости (SQL injection);
- оцени приоритет (blocker/major/minor);
- предложи минимально безопасное исправление
Формат ответа:
нумерованный список, для каждого бага (уязвимости) укажи: 
- blocker / major / minor
- что не так
- риск
- как исправить

Не предлагай крупный рефакторинг без причины, не придумывай новые библиотеки.

### Ответ ИИ

>Вот ревью кода на Go с точки зрения безопасности (OWASP Top 10).
>
>---
>
>### 1. A03:2021 – Injection (SQL Injection)
>- **Severity**: **blocker**
>- **Что не так**: Параметр `status` из URL-запроса напрямую вставляется в SQL-запрос через `fmt.Sprintf`. Злоумышленник может передать значение типа `"' OR '1'='1"` или `"'; DROP TABLE orders; --"`, что изменит логику запроса или выполнит произвольные SQL-команды.
>- **Риск**: Полный несанкционированный доступ к данным, обход аутентификации/авторизации, удаление или модификация таблиц, утечка всех заказов (включая чужие).
>- **Как исправить**:
>  - Использовать параметризованные запросы (плейсхолдеры `$1` для PostgreSQL, `?` для MySQL).
>  - Пример исправления:
>    ```go
>    func Orders(db *sql.DB, w http.ResponseWriter, r *http.Request) error {
>        status := r.URL.Query().Get("status")
>        query := "SELECT id, total FROM orders WHERE status = $1"
>        rows, err := db.Query(query, status)
>        if err != nil {
>            return err
>        }
>        defer rows.Close()
>        // обработка rows...
>        return nil
>    }
>    ```
>
>---
>
>### 2. A04:2021 – Insecure Design
>- **Severity**: **major**
>- **Что не так**: Отсутствует какая-либо проверка авторизации — кто может вызывать этот обработчик. Нет ограничений на допустимые значения `status` (например, только `pending`, `shipped`, `delivered`). Нет пагинации или лимита на количество возвращаемых записей.
>- **Риск**: Любой пользователь (включая неаутентифицированного) может получить доступ к заказам. Возможен отказ в обслуживании (DoS) через запрос, возвращающий миллионы строк.
>- **Как исправить**:
>  - Добавить middleware для проверки аутентификации и роли (например, только `admin` или `manager`).
>  - Валидировать `status` по белому списку допустимых значений.
>  - Добавить `LIMIT` и `OFFSET` для пагинации (или явно ограничить максимальное количество строк).
>
>---
>
>### 3. A09:2021 – Security Logging and Monitoring Failures
>- **Severity**: **minor**
>- **Что не так**: Нет логирования того, какой пользователь (IP, user ID) запросил заказы с каким параметром `status`. При успешной атаке или утечке невозможно определить источник и объём компрометации.
>- **Риск**: Затруднён анализ инцидентов, отсутствие возможности отследить злоумышленника, проблемы с соответствием регуляторным требованиям (GDPR, PCI DSS).
>- **Как исправить**:
>  - Добавить структурированное логирование (например, `slog.InfoContext`) с полями: `user_id` (из JWT/сессии), `remote_addr`, `status_param`, `rows_returned`, `duration_ms`.
>  - Логировать ошибки SQL с уровнем ERROR, но без вывода в ответ клиенту (чтобы не раскрывать структуру БД).
>
>---
>
>### Дополнительные замечания (вне OWASP, но важны):
>
>#### a) Ошибка обработки результата запроса
>- **Severity**: **major** (по качеству кода)
>- **Что не так**: Вызов `db.Query()` возвращает `*sql.Rows`, который не закрывается (`defer rows.Close()` отсутствует). Также результат запроса никак не используется — происходит утечка ресурсов (соединение с БД не освобождается).
>- **Риск**: Исчерпание пула соединений с БД, паника приложения, отказ в обслуживании после нескольких запросов.
>- **Как исправить**:
>  ```go
>  rows, err := db.Query(query, status)
>  if err != nil {
>      return err
>  }
>  defer rows.Close()
>  
>  var orders []Order
>  for rows.Next() {
>      var order Order
>      if err := rows.Scan(&order.ID, &order.Total); err != nil {
>          return err
>      }
>      orders = append(orders, order)
>  }
>  // записать orders в w (JSON, XML и т.д.)
>  ```
>
>#### b) Неверные кавычки в импорте
>- **Severity**: **minor**
>- **Что не так**: В коде используются типографские кавычки (`“` и `”`) вместо прямых двойных кавычек (`"`). Такой код не скомпилируется в Go.
>- **Риск**: Отсутствует (код просто не работает), но указывает на невнимательность разработчика.
>- **Как исправить**: Заменить все `“` и `”` на `"`.
>
>---
>
>### Итоговый приоритет:
>1. **SQL Injection** — blocker (исправить немедленно, до деплоя)
>2. **Отсутствие авторизации** — major
>3. **Утечка ресурсов (rows.Close)** — major
>4. **Логирование** — minor
>5. **Кавычки в импорте** — minor


### Карточка

Замечание ИИ | Тип           | «Светофор»   | Приоритет    | Что делаем
------------- | -------------| -------------| -------------| -------------
SQL Injection - параметр `status` из URL-запроса напрямую вставляется в SQL-запрос|корректное| accept| blocker| Использовать параметризованные запросы 
Отсутствие авторизации|корректное| accept| major| Добавить middleware для проверки аутентификации и роли, валидировать `status`
Отсутствие логирования|корректное| accept| minor| Добавить логирование
Утечка ресурсов (rows.Close)|корректное| accept| major| Добавить rows.Close()
Неверные кавычки в импорте|корректное| accept| minor| Исправить кавычки

### Комментарии

#### 1. Injection (SQL Injection)

- **Что не так**: Параметр `status` из URL-запроса напрямую вставляется в SQL-запрос через `fmt.Sprintf`. Злоумышленник может передать значение типа `"' OR '1'='1"` или `"'; DROP TABLE orders; --"`, что изменит логику запроса или выполнит произвольные SQL-команды.
- **Риск**: Полный несанкционированный доступ к данным, обход аутентификации/авторизации, удаление или модификация таблиц, утечка всех заказов (включая чужие).
- **Как исправить**: Использовать параметризованные запросы (плейсхолдеры `$1` для PostgreSQL, `?` для MySQL).
- **Приоритет**: blocker

#### 2. Insecure Design

- **Что не так**: Отсутствует какая-либо проверка авторизации — кто может вызывать этот обработчик. Нет ограничений на допустимые значения `status` (например, только `pending`, `shipped`, `delivered`). Нет пагинации или лимита на количество возвращаемых записей.
- **Риск**: Любой пользователь (включая неаутентифицированного) может получить доступ к заказам. Возможен отказ в обслуживании (DoS) через запрос, возвращающий миллионы строк.
- **Как исправить**:
  - Добавить middleware для проверки аутентификации и роли (например, только `admin` или `manager`).
  - Валидировать `status` по белому списку допустимых значений.
  - Добавить `LIMIT` и `OFFSET` для пагинации (или явно ограничить максимальное количество строк).
- **Приоритет**: major

#### 3. Security Logging and Monitoring Failures

- **Что не так**: Нет логирования того, какой пользователь (IP, user ID) запросил заказы с каким параметром `status`. При успешной атаке или утечке невозможно определить источник и объём компрометации.
- **Риск**: Затруднён анализ инцидентов, отсутствие возможности отследить злоумышленника, проблемы с соответствием регуляторным требованиям (GDPR, PCI DSS).
- **Как исправить**:
  - Добавить структурированное логирование (например, `slog.InfoContext`) с полями: `user_id` (из JWT/сессии), `remote_addr`, `status_param`, `rows_returned`, `duration_ms`.
  - Логировать ошибки SQL с уровнем ERROR, но без вывода в ответ клиенту (чтобы не раскрывать структуру БД).
- **Приоритет**: minor

#### 4. Ошибка обработки результата запроса

- **Что не так**: Вызов `db.Query()` возвращает `*sql.Rows`, который не закрывается (`defer rows.Close()` отсутствует). Также результат запроса никак не используется — происходит утечка ресурсов (соединение с БД не освобождается).
- **Риск**: Исчерпание пула соединений с БД, паника приложения, отказ в обслуживании после нескольких запросов.
- **Как исправить**: добавить `defer rows.Close()`
- **Приоритет**: minor

#### 5. Неверные кавычки в импорте

- **Что не так**: В коде используются типографские кавычки (`“` и `”`) вместо прямых двойных кавычек (`"`). Такой код не скомпилируется в Go.
- **Риск**: Отсутствует (код просто не работает), но указывает на невнимательность разработчика.
- **Как исправить**: Заменить все `“` и `”` на `"`.
- **Приоритет**: minor

### Next Steps

#### 1. Что исправляем ДО MERGE

##### Пункт 1: SQL Injection (Blocker)

Вместо
```go
query := fmt.Sprintf("SELECT * FROM orders WHERE status = '%s'", status)
rows, err := db.Query(query)
```
Сделать
```go
query := "SELECT * FROM orders WHERE status = $1"
rows, err := db.Query(query, status)
if err != nil {
    // обработка ошибки
}
```

##### Пункт 2: Insecure Design (Major)

Добавить проверку аутентификации и роли:

```go
func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("user_id")
    if userID == nil {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }
    
    role := r.Context().Value("role")
    if role != "admin" && role != "manager" {
        http.Error(w, "forbidden", http.StatusForbidden)
        return
    }
    
    // Валидация status по белому списку
    status := r.URL.Query().Get("status")
    allowedStatuses := map[string]bool{
        "pending":   true,
        "shipped":   true,
        "delivered": true,
    }
    if status != "" && !allowedStatuses[status] {
        http.Error(w, "invalid status value", http.StatusBadRequest)
        return
    }
    
    // Добавить LIMIT
    limit := 100
    query := "SELECT * FROM orders WHERE status = $1 LIMIT $2"
    rows, err := db.Query(query, status, limit)
}
```

##### Пункт 4: Утечка ресурсов (Minor)

```go
rows, err := db.Query(query, status)
if err != nil {
    return err
}
defer rows.Close()  // Добавить эту строку

for rows.Next() {
    // сканирование rows
}
```

##### Пункт 5: Неверные кавычки (Minor)

```go
// Вместо “text” и ”text”
// Использовать "text"

import "database/sql"  // Правильно
```

#### 2. Что выносим в ОТДЕЛЬНУЮ ЗАДАЧУ

##### Пункт 3: Security Logging (Minor)

```go
import "log/slog"

func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    
    // ... код обработчика ...
    
    // Логирование после выполнения
    slog.InfoContext(r.Context(), "orders fetched",
        "user_id", r.Context().Value("user_id"),
        "remote_addr", r.RemoteAddr,
        "status_param", r.URL.Query().Get("status"),
        "rows_returned", len(orders),
        "duration_ms", time.Since(start).Milliseconds(),
    )
    
    // Логирование ошибок SQL
    if err != nil {
        slog.ErrorContext(r.Context(), "sql query failed",
            "error", err,
            "user_id", r.Context().Value("user_id"),
        )
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }
}
```

#### 3. Какие тесты добавить

##### Тест на пункт 1 (SQL Injection)
```go
func TestGetOrdersHandler_SQLInjectionPrevention(t *testing.T) {
    req := httptest.NewRequest("GET", "/orders?status=pending' OR '1'='1", nil)
    rr := httptest.NewRecorder()
    
    handler := http.HandlerFunc(GetOrdersHandler)
    handler.ServeHTTP(rr, req)
    
    // Проверить, что запрос не привёл к ошибке БД или возврату всех записей
    // Проверить, что SQL-инъекция не сработала
    assert.Equal(t, http.StatusOK, rr.Code)
    // Убедиться, что вернулись только записи с ожидаемым статусом
}
```

##### Тест на пункт 2 (авторизация и валидация)
```go
func TestGetOrdersHandler_Unauthorized(t *testing.T) {
    req := httptest.NewRequest("GET", "/orders?status=pending", nil)
    rr := httptest.NewRecorder()
    
    // Без user_id в контексте
    GetOrdersHandler(rr, req)
    
    assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetOrdersHandler_InvalidRole(t *testing.T) {
    req := httptest.NewRequest("GET", "/orders?status=pending", nil)
    ctx := context.WithValue(req.Context(), "user_id", 123)
    ctx = context.WithValue(ctx, "role", "user")
    req = req.WithContext(ctx)
    rr := httptest.NewRecorder()
    
    GetOrdersHandler(rr, req)
    
    assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestGetOrdersHandler_InvalidStatus(t *testing.T) {
    req := httptest.NewRequest("GET", "/orders?status=invalid_status", nil)
    ctx := context.WithValue(req.Context(), "user_id", 123)
    ctx = context.WithValue(ctx, "role", "admin")
    req = req.WithContext(ctx)
    rr := httptest.NewRecorder()
    
    GetOrdersHandler(rr, req)
    
    assert.Equal(t, http.StatusBadRequest, rr.Code)
}
```

##### Тест на пункт 4 (ресурсы)
```go
func TestGetOrdersHandler_ClosesRows(t *testing.T) {
    // Использовать mock DB
    mock.ExpectQuery("SELECT .* FROM orders").
        WillReturnRows(sqlmock.NewRows([]string{"id", "status"}))
    
    // Вызвать обработчик
    
    // Проверить, что rows был закрыт
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

##### Тест на пункт 3 (логирование)
```go
func TestGetOrdersHandler_LogsAuditInfo(t *testing.T) {
    var logOutput bytes.Buffer
    logger := slog.New(slog.NewJSONHandler(&logOutput, nil))
    slog.SetDefault(logger)
    
    req := httptest.NewRequest("GET", "/orders?status=pending", nil)
    ctx := context.WithValue(req.Context(), "user_id", 123)
    ctx = context.WithValue(ctx, "role", "admin")
    req = req.WithContext(ctx)
    rr := httptest.NewRecorder()
    
    GetOrdersHandler(rr, req)
    
    // Проверить, что лог содержит audit-поля
    assert.Contains(t, logOutput.String(), "user_id")
    assert.Contains(t, logOutput.String(), "remote_addr")
    assert.Contains(t, logOutput.String(), "status_param")
}
```
