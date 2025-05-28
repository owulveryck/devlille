## 🧠 Présentation de MCP (Model Context Protocol) – 45 min devant \~100 développeurs

### 🎬 Introduction : une démo participative

Pour démarrer, je propose une démonstration interactive.

* J'affiche une page web simple contenant un **champ texte court**.
* Je demande à l’ensemble du public d’y saisir **leurs attentes vis-à-vis de cette présentation**.
* Toutes ces entrées sont centralisées dans une **base de données**.

Cette base est exposée via un protocole **JSON-RPC**, que je vais interroger en direct depuis mon terminal :

```bash
(
  cat <<\EOF
{"jsonrpc":"2.0","id":3,"method":"resources/read","params":{"name":"get_resources"}}
EOF
) | ./devlille
```

Ce retour contient toutes les attentes récoltées.

Je vais expliquer que cet appel est **l’équivalent d’un `GET`** dans une API REST, **mais en mode RPC**.

---

### 🔍 Zoom sur le protocole

* Le transport utilisé ici est **STDIO**, mais j’aurais aussi bien pu utiliser le **Web**.
* Le protocole d’appel est **JSON-RPC**.
* Ce que je montre ici est **la première brique de MCP : `resources/read`**.

#### Décorticage du message :

| Élément                      | JSON-RPC Standard | Spécifique MCP |
| ---------------------------- | ----------------- | -------------- |
| `"jsonrpc": "2.0"`           | ✅                 |                |
| `"id": 3`                    | ✅                 |                |
| `"method": "resources/read"` | ✅ (structure)     | ✅ (valeur)     |
| `"params": { "name": ... }`  | ✅ (structure)     | ✅ (valeurs)    |

---

### 🧠 L'IA entre en jeu

Je vais ensuite **demander à une IA en live** de :

1. Lire les attentes via MCP,
2. Les **synthétiser automatiquement**.

👉 On peut dire que **l’accès aux ressources MCP, c’est les yeux de l’IA** 😄

---

### 🛠 Deuxième action MCP : les *tools*

Je montre maintenant l’appel d’un **tool**, par exemple : modifier le slide en cours (je fournis l’ID du slide), en y injectant la synthèse des attentes.

> C’est **l’équivalent d’un `POST`** : on modifie l’environnement.
> Cela ouvre la voie à l’**agentivité**, c’est-à-dire la capacité d’agir pour une IA.

---

### 🗺 Troisième pilier MCP : le *prompt*

Le `prompt`, c’est **le mode d’emploi de l’agent**.

Je vais demander à mon assistant d’utiliser un prompt pour **générer une présentation complète** à partir des attentes.

> ✅ **Prompt à définir ici** – (je laisserai peut-être l’audience m’aider)

---

### 🚀 Conclusion

Le **MCP**, c’est :

* Une structuration simple : `resources`, `tools`, `prompts`
* Un protocole pensé pour l’**interaction entre IA et système**
* Une passerelle vers des usages nouveaux, tout comme les **API REST** l’ont été pour la transformation digitale

> On va créer des outils pensés non plus **pour les humains**, mais **pour les IA**.

---

### 🧠 Bonus : vers l’autonomie des agents

On n’a pas encore parlé des **notifications** envoyées par le serveur. Celles-ci permettent aux agents de **réagir de façon proactive** :

```go
for {
    Attendre une notification
    L’IA "réfléchit"
    tant que des choses sont à faire {
        L’IA agit (en appelant les bons serveurs MCP)
    }
}
```

> Cela permet à une IA de devenir **vraiment autonome**, comme un micro-agent intelligent.


