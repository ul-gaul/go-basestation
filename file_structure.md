# Structure des dossiers

```
go-basestation/
    > cmd/
    > config/
    > constants/
    > data/
        > collector/
        > packet/
        > parsing/
        > persistence/
    > pool/
    > resources/
        > samples/
    > ui/
        > plotting/
            > lines/
            > ticker/
        > views/
        > widgets/
    > utils/
```

## Contenu des dossiers
<dl>
<dt>cmd/</dt>
<dd>Interprète les lignes de commande.</dd>

<dt>config/</dt>
<dd>Lit les fichiers de configuration.</dd>

<dt>constants/</dt>
<dd>Contient les messages d'erreur et les constantes de l'application.</dd>

<dt>data/</dt>
<dd>
Contient tout ce qui est relié à l'obtention et au traitement des données.
<dl>
  <dt>data/collector/</dt>
  <dd>Contient le type Collector qui permet de regrouper et gérer les données reçues ou lues.</dd>

  <dt>data/packet/</dt>
  <dd>Contient les structures des packets de données reçues.</dd>

  <dt>data/parsing/</dt>
  <dd>Transforme les données binaires en structures Go et vice versa.</dd>

  <dt>data/persistence/</dt>
  <dd>Transmission des données et enregistrement dans des fichiers.</dd>
</dl>
</dd>

<dt>pool/</dt>
<dd>Contient les pools de threads. (multithreading)</dd>

<dt>resources/</dt>
<dd>Contient les ressources du projet (ex: données de tests, images, css, etc.)</dd>

<dt>ui/</dt>
<dd>Contient tout ce qui est relié à l'interface (frontend).</dd>
<dd><dl>
  <dt>ui/plotting/</dt>
  <dd>Adaptation de la library permettant de créer des plans cartésiens. (<a href="https://pkg.go.dev/gonum.org/v1/plot">gonum.org/v1/plot</a>)</dd>

  <dt>ui/views/</dt>
  <dd>Contient les éléments visuels spécifiques à l'application.</dd>

  <dt>ui/widgets/</dt>
  <dd>Contient des <a href="https://pkg.go.dev/gioui.org/widget"><em>widgets</em></a> pouvant être réutilisés.</dd>
</dl></dd>

<dt>utils/</dt>
<dd>Contient des fonctions utilitaires.</dd>

</dl>